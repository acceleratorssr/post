//go:build wireinject

package main

import (
	"github.com/google/wire"
	"post/user/grpc"
	"post/user/ioc"
	"post/user/repository"
	"post/user/repository/dao"
	"post/user/service"
)

func InitApp() *App {
	wire.Build(
		ioc.InitDB,
		ioc.InitGrpcServer,
		ioc.InitGrpcSSOClient,

		dao.NewUserGormDAO,
		repository.NewUserRepository,
		service.NewUserService,

		grpc.NewUserServiceServer,

		wire.Struct(new(App), "*"),
	)
	return new(App)
}
