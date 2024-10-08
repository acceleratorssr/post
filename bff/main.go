package main

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func main() {
	initViperWatch()
	app := InitApp()
	err := app.WebServer.Start()
	panic(err)
}

func initViperWatch() {
	cfile := pflag.String("config",
		"../bff/config/dev.yaml", "配置文件路径")
	pflag.Parse()

	viper.SetConfigFile(*cfile)
	viper.WatchConfig()
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}
