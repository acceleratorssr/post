//go:build wireinject

package main

import (
	"github.com/google/wire"
	"post/article/events"
	"post/article/grpc"
	"post/article/ioc"
	"post/article/repository"
	"post/article/repository/cache"
	"post/article/repository/dao"
	"post/article/service"
	distLock "post/pkg/redis_ex/distributed_lock"
)

var rankingServiceSet = wire.NewSet(
	cache.NewRankCache,
	cache.NewLocalCacheForRank,
	repository.NewBatchRankCache,
	service.NewBatchRankService)

var jobServiceSet = wire.NewSet(
	dao.NewGORMJobDAO,
	repository.NewPreemptJobRepository,
	service.NewCronJobService,
	ioc.InitLocalFuncExecutor,
	ioc.InitScheduler,

	ioc.InitRankingJob,
	ioc.InitJobs,
)

func InitApp() *App {
	wire.Build(
		distLock.NewClient,

		ioc.InitDB,
		ioc.InitKafka,
		ioc.NewKafkaSyncProducer,
		ioc.InitRedis,
		ioc.InitLikeClient,
		ioc.NewKafkaConsumer,
		events.NewKafkaReadProducer,
		events.NewKafkaPublishProducer,
		events.NewKafkaPublishedConsumer,

		rankingServiceSet,
		jobServiceSet,

		dao.NewSnowflakeNode0,
		dao.NewGORMArticleDao,
		cache.NewRedisArticleCache,
		repository.NewArticleAuthorRepository,
		repository.NewArticleReaderRepository,
		service.NewArticleService,

		grpc.NewArticleServiceServer,
		ioc.InitArticleService,

		wire.Struct(new(App), "*"),
	)
	return new(App)
}
