package main

import (
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
	"post/interactive/events"
	"post/internal/job"
	"post/internal/service"
	"post/internal/web"
)

type App struct {
	server         *gin.Engine
	consumers      []events.Consumer
	cron           *cron.Cron
	cronJobService service.JobService
	scheduler      *job.Scheduler

	// 测试用
	articleHandler *web.ArticleHandler
	db             *gorm.DB
}
