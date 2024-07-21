//go:build wireinject

package main

import (
	"github.com/google/wire"
	"post/events"
	"post/ioc"
	"post/repository"
	"post/repository/cache"
	"post/repository/dao"
	"post/service"
	"post/web"
)

func InitApp() *App {
	wire.Build(
		ioc.InitDB,
		//ioc.InitMongoDB,
		//ioc.InitLogger,
		ioc.InitKafka,
		ioc.InitRedis,
		ioc.NewKafkaSyncProducer,
		ioc.NewKafkaConsumer,
		ioc.InitWebServer,

		//events.NewKafkaConsumer,
		events.NewBatchKafkaConsumer,
		events.NewKafkaProducer,

		dao.NewGORMArticleDao,
		dao.NewGORMArticleLikeDao,
		//dao.NewMongoDB,
		cache.NewRedisArticleCache,
		//dao.NewSnowflakeNode,

		repository.NewArticleAuthorRepository,
		repository.NewArticleReaderRepository,
		repository.NewLikeRepository,

		service.NewArticleService,
		service.NewLikeService,

		web.NewArticleHandler,

		wire.Struct(new(App), "*"),
	)
	return new(App)
}
