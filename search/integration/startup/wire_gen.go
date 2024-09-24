// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package startup

import (
	"github.com/google/wire"
	"post/search/grpc"
	"post/search/ioc"
	"post/search/repository"
	"post/search/repository/dao"
	"post/search/service"
)

// Injectors from wire.go:

func InitSearchServer() *grpc.SearchServiceServer {
	client := InitESClient()
	articleDAO := dao.NewArticleElasticDAO(client)
	tagDAO := dao.NewTagESDAO(client)
	articleRepository := repository.NewArticleRepository(articleDAO, tagDAO)
	searchService := service.NewSearchService(articleRepository)
	searchServiceServer := grpc.NewSearchService(searchService)
	return searchServiceServer
}

func InitSyncServer() *grpc.SyncServiceServer {
	client := InitESClient()
	anyDAO := dao.NewAnyESDAO(client)
	anyRepository := repository.NewAnyRepository(anyDAO)
	articleDAO := dao.NewArticleElasticDAO(client)
	tagDAO := dao.NewTagESDAO(client)
	articleRepository := repository.NewArticleRepository(articleDAO, tagDAO)
	syncService := service.NewSyncService(anyRepository, articleRepository)
	syncServiceServer := grpc.NewSyncServiceServer(syncService)
	return syncServiceServer
}

// wire.go:

var serviceProviderSet = wire.NewSet(dao.NewArticleElasticDAO, dao.NewAnyESDAO, dao.NewTagESDAO, repository.NewAnyRepository, repository.NewArticleRepository, service.NewSyncService, service.NewSearchService)

var thirdProvider = wire.NewSet(
	InitESClient, ioc.InitLogger,
)
