// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"post/sso/config"
	"post/sso/grpc"
	"post/sso/ioc"
	"post/sso/repository"
	"post/sso/repository/dao"
	"post/sso/service"
)

// Injectors from wire.go:

func InitApp() *App {
	info := config.InitConfig()
	db := ioc.InitDB(info)
	ssoGormDAO := dao.NewSSOGormDAO(db)
	ssoRepository := repository.NewSSOGormRepository(ssoGormDAO)
	authUserService := service.NewAuthService(ssoRepository)
	authServiceServer := grpc.NewSSOServiceServer(authUserService, info)
	server := ioc.InitGrpcSSOServer(authServiceServer)
	app := &App{
		server: server,
	}
	return app
}
