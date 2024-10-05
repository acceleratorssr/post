package ioc

import (
	"context"
	_ "embed"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/opentelemetry/tracing"
	prom "gorm.io/plugin/prometheus"
	"log"
	"os"
	"post/article/repository/dao"
	"time"
)

//go:embed mysql.yaml
var mysqlDSN string

//go:embed mongoDB.yaml
var MongoDBDSN string

func InitDB() *gorm.DB {
	type Config struct {
		DSN string `yaml:"dsn"`
	}
	c := Config{
		DSN: mysqlDSN,
	}
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			LogLevel:                  logger.Info,
			SlowThreshold:             time.Second,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)
	db, err := gorm.Open(mysql.Open(c.DSN), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		panic(err)
	}

	err = dao.InitTables(db)
	if err != nil {
		panic(err)
	}

	// gorm自己的插件
	err = db.Use(prom.New(prom.Config{
		DBName:          "internal",
		RefreshInterval: 10,
		StartServer:     false, // 是否开启一个HTTP服务器来展示指标，main已经开了
		MetricsCollector: []prom.MetricsCollector{
			&prom.MySQL{
				VariableNames: []string{"Threads_connected", "Threads_running"},
			},
		},
	}))
	if err != nil {
		panic(err)
	}

	promCB := newCallbacks()
	prometheus.MustRegister(promCB.vector)

	// https://gorm.io/zh_CN/docs/write_plugins.html
	//监控查询执行时间
	err = db.Callback().Create().Before("*").
		Register("prometheus_create_before", promCB.before())
	err = db.Callback().Create().After("*").
		Register("prometheus_create_after", promCB.after("create"))
	err = db.Callback().Update().Before("*").
		Register("prometheus_update_before", promCB.before())
	err = db.Callback().Update().After("*").
		Register("prometheus_update_after", promCB.after("update"))
	err = db.Callback().Delete().Before("*").
		Register("prometheus_delete_before", promCB.before())
	err = db.Callback().Delete().After("*").
		Register("prometheus_delete_after", promCB.after("delete"))
	err = db.Callback().Raw().Before("*").
		Register("prometheus_raw_before", promCB.before())
	err = db.Callback().Raw().After("*").
		Register("prometheus_raw_after", promCB.after("raw"))
	err = db.Callback().Row().Before("*").
		Register("prometheus_row_before", promCB.before())
	err = db.Callback().Row().After("*").
		Register("prometheus_row_after", promCB.after("row"))

	if err != nil {
		panic(err)
	}

	err = db.Use(promCB)
	if err != nil {
		panic(err)
	}

	// https://github.com/go-gorm/opentelemetry
	// tracing.NewPlugin(tracing.WithDBName("internal")这个插件的实现同样使用底层的callback
	err = db.Use(tracing.NewPlugin(tracing.WithDBName("internal")))
	//tracing.WithoutQueryVariables(), // 不记录查询参数

	if err != nil {
		return nil
	}

	return db
}

func InitMongoDB() *mongo.Database {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	monitor := &event.CommandMonitor{
		// 命令执行前
		// 此处由于ctx没办法向下传递，所以不能用于传递开始时间
		Started: func(ctx context.Context, startedEvent *event.CommandStartedEvent) {
			fmt.Println(startedEvent.Command)
		},
		Succeeded: func(ctx context.Context, succeededEvent *event.CommandSucceededEvent) {

		},
		Failed: func(ctx context.Context, failedEvent *event.CommandFailedEvent) {

		},
	}
	opts := options.Client().ApplyURI(MongoDBDSN).SetMonitor(monitor)

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		panic(err)
	}

	db := client.Database("test")

	return db
}

// Callbacks 实现Plugin接口
type Callbacks struct {
	vector *prometheus.SummaryVec
}

func (c *Callbacks) Name() string {
	return "prometheus"
}

func (c *Callbacks) Initialize(db *gorm.DB) error {
	return nil
}

func newCallbacks() *Callbacks {
	vector := prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Subsystem: "http",
		Namespace: "post_service",
		Name:      "gorm_duration",
		Help:      "gorm_duration",
		Objectives: map[float64]float64{
			0.5:  0.05,
			0.9:  0.01,
			0.99: 0.001,
		},
	}, []string{"type", "table"})

	return &Callbacks{
		vector: vector,
	}
}
func (c *Callbacks) before() func(db *gorm.DB) {
	return func(db *gorm.DB) {
		start := time.Now()
		db.Set("startTime", start)
	}
}
func (c *Callbacks) after(typ string) func(db *gorm.DB) {
	return func(db *gorm.DB) {
		start, ok := db.Get("startTime")
		t := start.(time.Time)
		if !ok {
			return
		}

		table := db.Statement.Table
		if table == "" {
			table = "unknown"
		}
		c.vector.WithLabelValues(typ, table).
			Observe(float64(time.Since(t).Milliseconds()))
	}
}
