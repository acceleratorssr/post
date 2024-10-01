package ioc

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/viper"
	"post/bff/web"
	"post/pkg/gin_ex"
	"post/pkg/gin_ex/middleware"
	"post/pkg/logger"
	"time"
)

func InitGinServer(l logger.Logger, jwt *Jwt,
	user *web.UserHandler, sso *web.SSOHandler) *gin_ex.Server {
	engine := gin.Default()
	gin_ex.InitCounter(prometheus.CounterOpts{
		Namespace: "garden",
		Subsystem: "http_service",
		Name:      "http_request_count",
		Help:      "http_request_count",
	})

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	jwt.InitJwtValidateToken(ctx)
	cancel()

	engine.Use()

	jwtAOP := middleware.NewJwt(jwt.publicKey).Build()
	user.RegisterRoutes(ctx, engine)
	sso.RegisterRoutes(engine, jwtAOP)

	addr := viper.GetString("http.addr")

	return &gin_ex.Server{
		Engine: engine,
		Addr:   addr,
	}
}
