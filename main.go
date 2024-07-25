package main

import (
	"context"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"post/ioc"
	"time"
)

func main() {
	initPrometheus()
	fn := ioc.InitOTEL()

	app := InitApp()
	if app == nil {
		return
	}
	server := app.server
	topic := "article"
	for _, c := range app.consumers {
		err := c.Start(topic)
		// TODO 错误处理
		if err != nil {
			panic(err)
		}
	}
	err := server.Run(":9091")
	if err != nil {
		panic(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	fn(ctx)
}

func initPrometheus() {
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(":9191", nil)
	}()
}
