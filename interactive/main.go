package main

import (
	"fmt"
	"github.com/robfig/cron/v3"
	gin_extra "post/pkg/gin-extra"
	grpc_extra "post/pkg/grpc-extra"
	sarama_extra "post/pkg/sarama-extra"
)

func main() {
	app := InitApp()
	for _, c := range app.consumers {
		err := c.Start("")
		if err != nil {
			panic(err)
		}
	}

	app.cron.Start()

	go func() {
		fmt.Println("migrator start")
		app.webAdmin.Start()
	}()

	err := app.server.Serve()
	panic(err)
}

// App 控制main中的方法的启用，控制生命周期
type App struct {
	server    *grpc_extra.Server
	consumers []sarama_extra.Consumer
	webAdmin  *gin_extra.Server
	cron      *cron.Cron
}
