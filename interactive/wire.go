//go:build wireinject

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

var batchUpdateDBServiceSet = wire.NewSet(
	ioc.InitBatchUpdateDBJob,
	ioc.InitJobs,
)

var thirdPartySet = wire.NewSet(
	ioc.InitDoubleWritePool,
	ioc.InitDoubleWriteDB,
	ioc.InitBaseDB,
	ioc.InitTargetDB,

	ioc.InitGRPCexServer,
	ioc.InitRedis,
	ioc.InitLogger,

	events.NewKafkaIncrReadConsumer,
	ioc.NewKafkaConsumer,
	ioc.InitKafka,
	ioc.InitSyncProducer,
)

var likeSvcProvider = wire.NewSet(
	service.NewLikeService,
	repository.NewLikeRepository,
	dao.NewGORMArticleLikeDao,
	cache.NewRedisArticleLikeCache,
)

var migratorSet = wire.NewSet(
	ioc.InitMigratorServer,
	ioc.InitFixConsumer,
	ioc.InitMigratorProducer,
)

func InitApp() *App {
	wire.Build(
		thirdPartySet,
		likeSvcProvider,
		migratorSet,
		batchUpdateDBServiceSet,

		grpc.NewLikeServiceServer,

		wire.Struct(new(App), "*"),
	)
	return new(App)
}
