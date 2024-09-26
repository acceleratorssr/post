//go:build wireinject

package main

import (
	"github.com/google/wire"
	"post/sso/config"
	"post/sso/grpc"
	"post/sso/ioc"
	"post/sso/repository"
	"post/sso/repository/dao"
	"post/sso/service"
)

func InitApp() *App {
	wire.Build(
		config.InitConfig,
		service.NewAuthService,
		repository.NewSSOGormRepository,
		dao.NewSSOGormDAO,
		ioc.InitDB,
		ioc.InitGrpcSSOServer,
		grpc.NewSSOServiceServer,

		wire.Struct(new(App), "*"),
	)
	return new(App)
}
