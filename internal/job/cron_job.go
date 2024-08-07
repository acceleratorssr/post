package job

import "github.com/robfig/cron/v3"

type CronJobBuilder struct {
}

// NewCronJobBuilder 可加日志和监控
func NewCronJobBuilder() *CronJobBuilder {
	return &CronJobBuilder{}
}

// Build 构建一个cron.Job执行传入的job
func (C *CronJobBuilder) Build(job Job) cron.Job {
	return cron.FuncJob(func() {
		if err := job.Run(); err != nil {
			panic(err)
		}
	})
}
