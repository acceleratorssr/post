//go:build wireinject

package main

import (
	"github.com/google/wire"
	interRepo "post/interactive/repository"
	cache2 "post/interactive/repository/cache"
	interDAO "post/interactive/repository/dao"
	interService "post/interactive/service"
	"post/internal/events"
	"post/internal/ioc"
	distLock "post/internal/redis_distributed_lock"
	"post/internal/repository"
	"post/internal/repository/cache"
	"post/internal/repository/dao"
	"post/internal/service"
	"post/internal/web"
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
		//ioc.NewKafkaConsumer,
		ioc.InitWebServer,
		ioc.InitJobs,
		ioc.InitRankingJob,
		ioc.InitScheduler,
		ioc.InitLocalFuncExecutor,
		ioc.InitIntrGRPCClient,

		rankingServiceSet,

		//events.NewKafkaConsumer,
		//events2.NewBatchKafkaConsumer,
		events.NewKafkaProducer,

		dao.NewGORMArticleDao,
		interDAO.NewGORMArticleLikeDao,
		dao.NewGORMJobDAO,
		//dao.NewMongoDB,
		cache.NewRedisArticleCache,
		cache.NewLocalCacheForRank,
		cache2.NewRedisArticleLikeCache,
		//dao.NewSnowflakeNode,

		repository.NewArticleAuthorRepository,
		repository.NewArticleReaderRepository,
		interRepo.NewLikeRepository,
		repository.NewPreemptJobRepository,

		service.NewArticleService,
		interService.NewLikeService,
		service.NewCronJobService,

		web.NewArticleHandler,

		wire.Struct(new(App), "*"),
	)
	return new(App)
}
