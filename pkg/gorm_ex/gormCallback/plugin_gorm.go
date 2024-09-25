package gormex

import (
	_ "embed"
	"github.com/prometheus/client_golang/prometheus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

func InitDB(mysqlDSN string) *gorm.DB {
	db, err := gorm.Open(mysql.Open(mysqlDSN), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	promCB := NewCallbacks()
	prometheus.MustRegister(promCB.vector)

	err = db.Callback().Query().Before("*").
		Register("prometheus_query_before", promCB.before())
	err = db.Callback().Query().After("*").
		Register("prometheus_query_after", promCB.after("query"))
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

	return db
}

type Option func(callbacks *Callbacks)

// Callbacks 实现Plugin接口
type Callbacks struct {
	vector      *prometheus.SummaryVec
	subsystem   string
	namespace   string
	name        string
	help        string
	constLabels map[string]string
	objectives  map[float64]float64
}

func (c *Callbacks) Name() string {
	return "prometheus"
}

func (c *Callbacks) Initialize(db *gorm.DB) error {
	return nil
}

func NewCallbacks(opts ...Option) *Callbacks {
	vector := prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Subsystem: "http",
		Namespace: "service",
		Name:      "gorm_duration",
		Help:      "gorm_duration",
		ConstLabels: map[string]string{
			"instance_id": "internal",
		},
		Objectives: map[float64]float64{
			0.5:   0.01,
			0.75:  0.01,
			0.90:  0.005,
			0.98:  0.002,
			0.99:  0.001,
			0.999: 0.0001,
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

func WithSubsystem(subsystem string) Option {
	return func(c *Callbacks) {
		c.subsystem = subsystem
	}
}
func WithNamespace(namespace string) Option {
	return func(c *Callbacks) {
		c.namespace = namespace
	}
}

func WithName(name string) Option {
	return func(c *Callbacks) {
		c.name = name
	}
}

func WithHelp(help string) Option {
	return func(c *Callbacks) {
		c.help = help
	}
}

func WithConstLabels(constLabels map[string]string) Option {
	return func(c *Callbacks) {
		c.constLabels = constLabels
	}
}

func WithObjectives(objectives map[float64]float64) Option {
	return func(c *Callbacks) {
		c.objectives = objectives
	}
}
