package ioc

import (
	_ "embed"
	"github.com/prometheus/client_golang/prometheus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/plugin/opentelemetry/tracing"
	prom "gorm.io/plugin/prometheus"
	dao2 "post/interactive/repository/dao"
	"post/pkg/gorm-extra/connpool"
	"time"
)

//go:embed mysql.yaml
var mysqlDSN string

//go:embed target_mysql.yaml
var target string

// TargetDB 为了区分两个不同的DB
type TargetDB *gorm.DB

type BaseDB *gorm.DB

func InitDoubleWriteDB(pool *connpool.DoubleWritePool) *gorm.DB {
	db, err := gorm.Open(mysql.New(mysql.Config{
		Conn: pool,
	}))
	if err != nil {
		panic(err)
	}
	return db
}

func InitTargetDB() TargetDB {
	db, err := gorm.Open(mysql.Open(target), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	err = dao2.InitTables(db)
	return db
}

func InitBaseDB() BaseDB {
	//type Config struct {
	//	DSN string `yaml:"dsn"`
	//}
	//c := Config{
	//	DSN: mysqlDSN,
	//}

	db, err := gorm.Open(mysql.Open(mysqlDSN), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	err = dao2.InitTables(db)
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
