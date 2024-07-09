package ioc

import (
	"github.com/gin-gonic/gin"
	"post/web"
)

func InitWebServer(middlewares []gin.HandlerFunc, handler *web.ArticleHandler) *gin.Engine {
	server := gin.Default()
	server.Use(middlewares...)
	handler.RegisterRoutes(server)

	return server
}
