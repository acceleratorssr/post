package ioc

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"post/pkg/logger"
)

func InitLogger() logger.Logger {
	cfg := zap.NewProductionConfig() // 生产模式，默认输出为 JSON 格式

	// todo k8s 会从标准输出中采集日志，更改部署后需修改此处
	//cfg.OutputPaths = []string{"/var/log/bff.log"}

	err := viper.UnmarshalKey("log", &cfg)
	if err != nil {
		panic(err)
	}

	l, err := cfg.Build()
	if err != nil {
		panic(err)
	}

	return logger.NewZapLogger(l)
}
