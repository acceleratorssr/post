package main

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"os"
	"post/pkg/grpc-extra"
	"post/search/events"
	"strings"
)

func main() {
	initViperWatch()
	app := Init()
	for _, c := range app.consumers {
		err := c.Start()
		if err != nil {
			panic(err)
		}
	}
	err := app.server.Serve()
	panic(err)
}

func initViperWatch() {
	execPath, err := os.Executable()
	if err != nil {
		panic(err)
	}

	searchPath := ""
	if !strings.Contains(execPath, "output") {
		searchPath = "./search/config/dev.yaml"
	} else {
		searchPath = "../search/config/dev.yaml"
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

type App struct {
	server    *grpc_extra.Server
	consumers []events.Consumer
}
