// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

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

// Injectors from wire.go:

func Init() *App {
	client := ioc.InitESClient()
	anyDAO := dao.NewAnyESDAO(client)
	anyRepository := repository.NewAnyRepository(anyDAO)
	articleDAO := dao.NewArticleElasticDAO(client)
	articleRepository := repository.NewArticleRepository(articleDAO)
	syncService := service.NewSyncService(anyRepository, articleRepository)
	syncServiceServer := grpc.NewSyncServiceServer(syncService)
	searchService := service.NewSearchService(articleRepository)
	searchServiceServer := grpc.NewSearchService(searchService)
	server := ioc.InitGRPCxServer(syncServiceServer, searchServiceServer)
	saramaClient := ioc.InitKafka()
	logger := ioc.InitLogger()
	articleConsumer := events.NewArticleConsumer(saramaClient, logger, syncService)
	v := ioc.NewConsumers(articleConsumer)
	app := &App{
		server:    server,
		consumers: v,
	}
	return app
}

// wire.go:

var serviceProviderSet = wire.NewSet(dao.NewArticleElasticDAO, dao.NewAnyESDAO, repository.NewArticleRepository, repository.NewAnyRepository, service.NewSyncService, service.NewSearchService)

var thirdProvider = wire.NewSet(ioc.InitESClient, ioc.InitLogger, ioc.InitKafka)