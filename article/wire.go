//go:build wireinject

package main

import (
	"github.com/google/wire"
	"post/article/events"
	"post/article/grpc"
	"post/article/ioc"
	"post/article/repository"
	"post/article/repository/cache"
	"post/article/repository/cache/compression"
	"post/article/repository/dao"
	"post/article/service"
	distLock "post/pkg/redis-extra/distributed_lock"
)

var rankingServiceSet = wire.NewSet(
	cache.NewRankCache,
	cache.NewLocalCacheForRank,
	repository.NewBatchRankCache,
	service.NewBatchRankService,
)

var schedulerServiceSet = wire.NewSet(
	dao.NewGORMJobDAO,
	repository.NewPreemptJobRepository,
	service.NewCronJobService,
	ioc.InitLocalFuncExecutor,
	ioc.InitScheduler,
)

var jobServiceSet = wire.NewSet(
	ioc.InitRankingJob,
	ioc.InitJobs,
)

var smallMessagesSet = wire.NewSet(
	events.NewKafkaSyncProducerForSmallMessages,
	events.NewKafkaReadProducer,
	events.NewKafkaRecommendProducer,
)

var largeMessagesSet = wire.NewSet(
	events.NewKafkaSyncProducerForLargeMessages,
	events.NewKafkaPublishProducer,
)

func InitApp() *App {
	wire.Build(
		distLock.NewClient,

		ioc.InitDB,
		ioc.InitKafka,
		ioc.InitRedis,
		ioc.InitLikeClient,

		events.NewKafkaPublishedConsumer,
		events.NewKafkaConsumer,

		rankingServiceSet,

		schedulerServiceSet,
		jobServiceSet,

		smallMessagesSet,
		largeMessagesSet,

		dao.NewSnowflakeNode0,
		dao.NewGORMArticleDao,
		compression.NewArticleCompressionByGZIP,
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
