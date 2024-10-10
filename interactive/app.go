package main

import (
	"github.com/robfig/cron/v3"
	"post/pkg/gin-extra"
	"post/pkg/grpc-extra"
	"post/pkg/sarama-extra"
)

// App 控制main中的方法的启用，控制生命周期
type App struct {
	server    *grpc_extra.Server
	consumers []sarama_extra.Consumer
	webAdmin  *gin_extra.Server
	cron      *cron.Cron
}
