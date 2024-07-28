package main

import (
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	"post/events"
)

type App struct {
	server    *gin.Engine
	consumers []events.Consumer
	cron      *cron.Cron
}
