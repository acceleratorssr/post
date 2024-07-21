package main

import (
	"github.com/gin-gonic/gin"
	"post/events"
)

type App struct {
	server    *gin.Engine
	consumers []events.Consumer
}
