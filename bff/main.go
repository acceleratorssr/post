package main

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"os"
	"strings"
)

func main() {
	initViperWatch()
	app := InitApp()
	err := app.WebServer.Start()
	panic(err)
}

func initViperWatch() {
	execPath, err := os.Executable()
	if err != nil {
		panic(err)
	}

	searchPath := ""
	if !strings.Contains(execPath, "output") {
		searchPath = "./bff/config/dev.yaml"
	} else {
		searchPath = "../bff/config/dev.yaml"
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
