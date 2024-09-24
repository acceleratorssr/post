//go:build wireinject

package startup

import (
	"github.com/google/wire"
	"post/search/grpc"
	"post/search/ioc"
	"post/search/repository"
	"post/search/repository/dao"
	"post/search/service"
)

var serviceProviderSet = wire.NewSet(
	dao.NewArticleElasticDAO,
	dao.NewAnyESDAO,
	dao.NewTagESDAO,
	repository.NewAnyRepository,
	repository.NewArticleRepository,
	service.NewSyncService,
	service.NewSearchService,
)

var thirdProvider = wire.NewSet(
	InitESClient,
	ioc.InitLogger)

func InitSearchServer() *grpc.SearchServiceServer {
	wire.Build(
		thirdProvider,
		serviceProviderSet,
		grpc.NewSearchService,
	)
	return new(grpc.SearchServiceServer)
}

func InitSyncServer() *grpc.SyncServiceServer {
	wire.Build(
		thirdProvider,
		serviceProviderSet,
		grpc.NewSyncServiceServer,
	)
	return new(grpc.SyncServiceServer)
}
