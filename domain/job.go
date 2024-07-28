package domain

import (
	"github.com/robfig/cron/v3"
	"time"
)

type Job struct {
	ID       int64
	Name     string
	Executor string
	Cron     string

	Cfg         string
	Status      int
	ExecuteTime int64
	Version     int64

	CancelFunc func() error

	Ctime int64
	Utime int64
}

func (j *Job) NextExecTime() time.Time {
	parser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)
	s, _ := parser.Parse(j.Cron)
	return s.Next(time.Now())
}
