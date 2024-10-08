package main

import (
	"github.com/robfig/cron/v3"
	"post/article/job"
	"post/article/service"
	"post/pkg/grpc-extra"
	"post/pkg/sarama-extra"
)

type App struct {
	server         *grpc_extra.Server
	consumers      []sarama_extra.Consumer
	cron           *cron.Cron
	cronJobService service.JobService
	scheduler      *job.Scheduler
}
