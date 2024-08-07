//go:build wireinject

package startup

import (
	"github.com/google/wire"
	"post/interactive/grpc"
	"post/interactive/repository"
	"post/interactive/repository/cache"
	"post/interactive/repository/dao"
	"post/interactive/service"
)

var thirdProvider = wire.NewSet(
	InitRedis, InitDB,
	InitLogger,
	InitKafka,
)

func InitGRPCServer() *grpc.InteractiveServiceServer {
	wire.Build(
		grpc.NewInteractiveServiceServer,
		thirdProvider,
		dao.NewGORMInteractiveDAO,
		cache.NewRedisInteractiveCache,
		repository.NewCachedInteractiveRepository,
		service.NewInteractiveService,
	)
	return new(grpc.InteractiveServiceServer)
}
