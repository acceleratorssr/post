package ioc

import (
	"go.uber.org/zap"
	"post/logger"
)

func InitLogger() logger.Logger {
	l, _ := zap.NewProduction()
	//defer l.Sync()

	return logger.NewZapLogger(l)
}
