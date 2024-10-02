package main

import (
	"context"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"post/article/ioc"
	"time"
)

func main() {
	initPrometheus()
	fn := ioc.InitOTEL()

	app := InitApp()
	if app == nil {
		return
	}

	app.cron.Start()

	err := app.server.Serve()
	if err != nil {
		panic(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	fn(ctx)
	<-app.cron.Stop().Done()
}

func initPrometheus() {
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(":9191", nil)
	}()
}
