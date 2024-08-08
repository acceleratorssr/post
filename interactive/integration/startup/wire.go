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

func InitGRPCServer() *grpc.LikeServiceServer {
	wire.Build(
		grpc.NewLikeServiceServer,
		thirdProvider,
		dao.NewGORMArticleLikeDao,
		cache.NewRedisArticleLikeCache,
		repository.NewLikeRepository,
		service.NewLikeService,
	)
	return new(grpc.LikeServiceServer)
}
