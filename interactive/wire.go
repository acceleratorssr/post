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

var thirdPartySet = wire.NewSet(
	ioc.InitDB,
	ioc.InitRedis,
	ioc.InitLogger,
	ioc.InitKafka,
	ioc.InitGRPCexServer,
	ioc.NewKafkaConsumer,
)

var likeSvcProvider = wire.NewSet(
	service.NewLikeService,
	repository.NewLikeRepository,
	dao.NewGORMArticleLikeDao,
	cache.NewRedisArticleLikeCache,
)

func InitApp() *App {
	wire.Build(
		likeSvcProvider,
		thirdPartySet,
		events.NewKafkaConsumer,

		grpc.NewLikeServiceServer,
		wire.Struct(new(App), "*"),
	)
	return new(App)
}
