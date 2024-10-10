// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/google/wire"
	"post/interactive/events"
	"post/interactive/grpc"
	"post/interactive/ioc"
	"post/interactive/repository"
	"post/interactive/repository/cache"
	"post/interactive/repository/dao"
	"post/interactive/service"
)

// Injectors from wire.go:

func InitApp() *App {
	baseDB := ioc.InitBaseDB()
	targetDB := ioc.InitTargetDB()
	doubleWritePool := ioc.InitDoubleWritePool(baseDB, targetDB)
	db := ioc.InitDoubleWriteDB(doubleWritePool)
	articleLikeDao := dao.NewGORMArticleLikeDao(db)
	cmdable := ioc.InitRedis()
	articleLikeCache := cache.NewRedisArticleLikeCache(cmdable)
	likeRepository := repository.NewLikeRepository(articleLikeDao, articleLikeCache)
	likeService := service.NewLikeService(likeRepository)
	likeServiceServer := grpc.NewLikeServiceServer(likeService)
	server := ioc.InitGRPCexServer(likeServiceServer)
	client := ioc.InitKafka()
	kafkaReadConsumer := events.NewKafkaIncrReadConsumer(client, likeRepository)
	consumer := ioc.InitFixConsumer(baseDB, targetDB, client)
	v := ioc.NewKafkaConsumer(kafkaReadConsumer, consumer)
	syncProducer := ioc.InitSyncProducer(client)
	inconsistentProducer := ioc.InitMigratorProducer(syncProducer)
	gin_extraServer := ioc.InitMigratorServer(baseDB, targetDB, doubleWritePool, inconsistentProducer)
	batchUpdateDBJob := ioc.InitBatchUpdateDBJob(likeService, articleLikeCache)
	cron := ioc.InitJobs(batchUpdateDBJob)
	app := &App{
		server:    server,
		consumers: v,
		webAdmin:  gin_extraServer,
		cron:      cron,
	}
	return app
}

// wire.go:

var batchUpdateDBServiceSet = wire.NewSet(ioc.InitBatchUpdateDBJob, ioc.InitJobs)

var thirdPartySet = wire.NewSet(ioc.InitDoubleWritePool, ioc.InitDoubleWriteDB, ioc.InitBaseDB, ioc.InitTargetDB, ioc.InitGRPCexServer, ioc.InitRedis, ioc.InitLogger, events.NewKafkaIncrReadConsumer, ioc.NewKafkaConsumer, ioc.InitKafka, ioc.InitSyncProducer)

var likeSvcProvider = wire.NewSet(service.NewLikeService, repository.NewLikeRepository, dao.NewGORMArticleLikeDao, cache.NewRedisArticleLikeCache)

var migratorSet = wire.NewSet(ioc.InitMigratorServer, ioc.InitFixConsumer, ioc.InitMigratorProducer)
