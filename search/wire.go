//go:build wireinject

package main

import (
	"github.com/google/wire"
	"post/search/events"
	"post/search/grpc"
	"post/search/ioc"
	"post/search/repository"
	"post/search/repository/dao"
	"post/search/service"
)

var serviceProviderSet = wire.NewSet(
	dao.NewArticleElasticDAO,
	dao.NewAnyESDAO,
	repository.NewArticleRepository,
	repository.NewAnyRepository,
	service.NewSyncService,
	service.NewSearchService,
)

var thirdProvider = wire.NewSet(
	ioc.InitESClient,
	ioc.InitLogger,
	ioc.InitKafka)

func Init() *App {
	wire.Build(
		thirdProvider,
		serviceProviderSet,
		grpc.NewSyncServiceServer,
		grpc.NewSearchService,
		events.NewArticleConsumer,
		ioc.InitGRPCxServer,
		ioc.NewConsumers,
		wire.Struct(new(App), "*"),
	)
	return new(App)
}
