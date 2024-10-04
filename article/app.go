package main

import (
	"github.com/robfig/cron/v3"
	"post/article/job"
	"post/article/service"
	"post/pkg/grpc_ex"
	"post/pkg/sarama_ex"
)

type App struct {
	server         *grpc_ex.Server
	consumers      []sarama_ex.Consumer
	cron           *cron.Cron
	cronJobService service.JobService
	scheduler      *job.Scheduler
}
