//go:build wireinject

package main

import (
	"github.com/google/wire"
	"post/sso/config"
	"post/sso/grpc"
	"post/sso/ioc"
	"post/sso/repository"
	"post/sso/repository/cache"
	"post/sso/repository/dao"
	"post/sso/service"
)

func InitApp() *App {
	wire.Build(
		config.InitConfig,
		service.NewAuthService,
		service.NewJWTService,

		repository.NewSSOGormRepository,
		repository.NewSSOCache,

		cache.NewRedisCache,
		dao.NewSSOGormDAO,

		ioc.InitDB,
		ioc.InitRedis,
		ioc.InitGrpcSSOServer,

		grpc.NewSSOServiceServer,

		wire.Struct(new(App), "*"),
	)
	return new(App)
}
