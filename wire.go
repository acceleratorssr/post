//go:build wireinject

package main

import (
	"github.com/google/wire"
	"post/events"
	"post/ioc"
	distLock "post/redis_distributed_lock"
	"post/repository"
	"post/repository/cache"
	"post/repository/dao"
	"post/service"
	"post/web"
)

var rankingServiceSet = wire.NewSet(
	cache.NewRankCache,
	repository.NewBatchRankCache,
	service.NewBatchRankService)

func InitApp() *App {
	wire.Build(
		distLock.NewClient,

		ioc.InitDB,
		//ioc.InitMongoDB,
		//ioc.InitLogger,
		ioc.InitKafka,
		ioc.InitRedis,
		ioc.NewKafkaSyncProducer,
		ioc.NewKafkaConsumer,
		ioc.InitWebServer,
		ioc.InitJobs,
		ioc.InitRankingJob,

		rankingServiceSet,

		//events.NewKafkaConsumer,
		events.NewBatchKafkaConsumer,
		events.NewKafkaProducer,

		dao.NewGORMArticleDao,
		dao.NewGORMArticleLikeDao,
		//dao.NewMongoDB,
		cache.NewRedisArticleCache,
		cache.NewLocalCacheForRank,
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
