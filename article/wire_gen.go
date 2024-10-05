// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

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
	"post/pkg/redis_ex/distributed_lock"
)

// Injectors from wire.go:

func InitApp() *App {
	db := ioc.InitDB()
	articleDao := dao.NewGORMArticleDao(db)
	cmdable := ioc.InitRedis()
	articleCache := cache.NewRedisArticleCache(cmdable)
	node := dao.NewSnowflakeNode0()
	articleAuthorRepository := repository.NewArticleAuthorRepository(articleDao, articleCache, node)
	articleReaderRepository := repository.NewArticleReaderRepository(articleDao)
	client := ioc.InitKafka()
	smallMessagesProducer := ioc.NewKafkaSyncProducerForSmallMessages(client)
	readProducer := events.NewKafkaReadProducer(smallMessagesProducer)
	largeMessagesProducer := ioc.NewKafkaSyncProducerForLargeMessages(client)
	publishedProducer := events.NewKafkaPublishProducer(largeMessagesProducer)
	articleService := service.NewArticleService(articleAuthorRepository, articleReaderRepository, readProducer, publishedProducer)
	articleServiceServer := grpc.NewArticleServiceServer(articleService)
	server := ioc.InitArticleService(articleServiceServer)
	kafkaPublishedConsumer := events.NewKafkaPublishedConsumer(client, articleReaderRepository)
	v := events.NewKafkaConsumer(kafkaPublishedConsumer)
	likeServiceClient := ioc.InitLikeClient()
	rankCache := cache.NewRankCache(cmdable)
	localCacheForRank := cache.NewLocalCacheForRank()
	rankRepository := repository.NewBatchRankCache(rankCache, localCacheForRank)
	rankService := service.NewBatchRankService(articleService, likeServiceClient, rankRepository)
	distributed_lockClient := distributed_lock.NewClient(cmdable)
	rankingJob := ioc.InitRankingJob(rankService, distributed_lockClient)
	cron := ioc.InitJobs(rankingJob)
	jobDAO := dao.NewGORMJobDAO(db)
	jobRepository := repository.NewPreemptJobRepository(jobDAO)
	jobService := service.NewCronJobService(jobRepository)
	localFuncExecutor := ioc.InitLocalFuncExecutor(rankService)
	scheduler := ioc.InitScheduler(jobService, localFuncExecutor)
	app := &App{
		server:         server,
		consumers:      v,
		cron:           cron,
		cronJobService: jobService,
		scheduler:      scheduler,
	}
	return app
}

// wire.go:

var rankingServiceSet = wire.NewSet(cache.NewRankCache, cache.NewLocalCacheForRank, repository.NewBatchRankCache, service.NewBatchRankService)

var schedulerServiceSet = wire.NewSet(dao.NewGORMJobDAO, repository.NewPreemptJobRepository, service.NewCronJobService, ioc.InitLocalFuncExecutor, ioc.InitScheduler)

var jobServiceSet = wire.NewSet(ioc.InitRankingJob, ioc.InitJobs)

var smallMessagesSet = wire.NewSet(ioc.NewKafkaSyncProducerForSmallMessages, events.NewKafkaReadProducer)

var largeMessagesSet = wire.NewSet(ioc.NewKafkaSyncProducerForLargeMessages, events.NewKafkaPublishProducer)
