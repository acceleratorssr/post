package ioc

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"post/pkg/gin_ex"
	"post/pkg/metric"
	"post/web"
	"time"
)

func InitWebServer(articleHdl *web.ArticleHandler) *gin.Engine {
	server := gin.Default()
	//mdls []gin.HandlerFunc,
	gin_ex.InitCounter(prometheus.CounterOpts{
		Name:      "http_request_count",
		Help:      "http_request_count",
		Subsystem: "http",
		Namespace: "post_service",
	})

	mdls := []gin.HandlerFunc{
		corsHdl(),
		// https://github.com/open-telemetry/opentelemetry-go-contrib/blob/main/instrumentation/github.com/gin-gonic/gin/otelgin/gintrace.go
		otelgin.Middleware("post"),
		(&metric.Metric{
			Subsystem:  "http",
			Namespace:  "post_service",
			Name:       "http_request_duration_seconds",
			Help:       "http_request_duration",
			InstanceID: "post",
		}).Build(),
	}
	server.Use(mdls...)
	articleHdl.RegisterRoutes(server)
	return server
}
func corsHdl() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"token"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,

		AllowOrigins: []string{"http://127.0.0.1"},
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
	})
}

func InitMiddleware() []gin.HandlerFunc {
	return []gin.HandlerFunc{
		corsHdl(),
		(&metric.Metric{
			Name:       "http_request_duration_seconds",
			Help:       "http request duration",
			Subsystem:  "http",
			Namespace:  "post_service",
			InstanceID: "post",
		}).Build(),
	}
}
