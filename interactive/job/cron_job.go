package job

import "github.com/robfig/cron/v3"

type CronJobBuilder struct {
}

func NewCronJobBuilder() *CronJobBuilder {
	return &CronJobBuilder{}
}

func (C *CronJobBuilder) Build(job Job) cron.Job {
	return cron.FuncJob(func() {
		if err := job.Run(); err != nil {
			panic(err)
		}
	})
}
