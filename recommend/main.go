package main

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"os"
	grpcextra "post/pkg/grpc-extra"
	saramaextra "post/pkg/sarama-extra"
	"strings"
)

type App struct {
	server    *grpcextra.Server
	consumers []saramaextra.Consumer
}

func main() {
	initViperWatch()
	app := InitApp()
	if app == nil {
		return
	}

	for _, c := range app.consumers {
		err := c.Start("")
		if err != nil {
			panic(err)
		}
	}

	err := app.server.Serve()
	if err != nil {
		panic(err)
	}
}

func initViperWatch() {
	execPath, err := os.Executable()
	if err != nil {
		panic(err)
	}

	searchPath := ""
	if !strings.Contains(execPath, "output") {
		searchPath = "./recommend/config/dev.yaml"
	} else {
		searchPath = "../recommend/config/dev.yaml"
	}

	cfile := pflag.String("config", searchPath, "yaml path")
	pflag.Parse()
	viper.SetConfigFile(*cfile)
	viper.WatchConfig()
	err = viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}
