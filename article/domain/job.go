package domain

import (
	"github.com/robfig/cron/v3"
	"time"
)

type Job struct {
	ID       int64
	Name     string
	Executor string
	Cron     string // 存储cron表达式，即方便执行完后自动设置下一次执行的时间

	Cfg         string
	Status      int
	ExecuteTime int64
	Version     int64

	Topic     string
	Partition string
	Data      string

	CancelFunc func() error

	Ctime int64
	Utime int64
}

func (j *Job) NextExecTime() time.Time {
	parser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)
	s, _ := parser.Parse(j.Cron)
	return s.Next(time.Now())
}
