package ioc

import (
	"context"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/viper"
	"post/bff/web"
	"post/pkg/gin-extra"
	"post/pkg/gin-extra/middleware"
	"post/pkg/logger"
	"time"
)

func InitGinServer(l logger.Logger, jwt *Jwt,
	user *web.UserHandler, sso *web.SSOHandler,
	article *web.ArticleHandler) *gin_extra.Server {

	engine := gin.Default()
	gin_extra.InitCounter(prometheus.CounterOpts{
		Namespace: "garden",
		Subsystem: "http_service",
		Name:      "http_request_count",
		Help:      "http_request_count",
	})

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	jwt.InitJwtValidateToken(ctx)
	cancel()

	mw := []gin.HandlerFunc{
		corsHdl(),
	}
	engine.Use(mw...)

	jwtAOP := middleware.NewJwt(jwt.publicKey).Build()
	user.RegisterRoutes(engine, jwtAOP)
	sso.RegisterRoutes(engine, jwtAOP)
	article.RegisterRoutes(engine, jwtAOP)

	addr := viper.GetString("http.addr")

	return &gin_extra.Server{
		Engine: engine,
		Addr:   addr,
	}
}

func corsHdl() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Authorization"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,

		AllowOrigins: []string{"http://127.0.0.1"},
		AllowMethods: []string{"GET", "POST"},
	})
}
