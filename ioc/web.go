package ioc

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"post/pkg/metric"
	"post/web"
	"time"
)

func InitWebServer(articleHdl *web.ArticleHandler) *gin.Engine {
	server := gin.Default()
	//mdls []gin.HandlerFunc,
	mdls := []gin.HandlerFunc{
		corsHdl(),
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
