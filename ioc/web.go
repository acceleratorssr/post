package ioc

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"post/web"
	"time"
)

func InitWebServer(articleHdl *web.ArticleHandler) *gin.Engine {
	server := gin.Default()
	articleHdl.RegisterRoutes(server)
	return server
}
func corsHdl() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		ExposeHeaders:    []string{"token"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})
}
