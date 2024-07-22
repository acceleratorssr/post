package main

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

func main() {
	initPrometheus()
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
}

func initPrometheus() {
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(":9191", nil)
	}()
}
