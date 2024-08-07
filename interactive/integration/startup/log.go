package startup

import (
	"go.uber.org/zap"
	"post/pkg/logger"
)

func InitLogger() logger.Logger {
	l, _ := zap.NewProduction()

	return logger.NewZapLogger(l)
}
