// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"post/user/grpc"
	"post/user/ioc"
	"post/user/repository"
	"post/user/repository/dao"
	"post/user/service"
)

// Injectors from wire.go:

func InitApp() *App {
	db := ioc.InitDB()
	userGormDAO := dao.NewUserGormDAO(db)
	userRepository := repository.NewUserRepository(userGormDAO)
	userService := service.NewUserService(userRepository)
	authServiceClient := ioc.InitGrpcSSOClient()
	userServiceServer := grpc.NewUserServiceServer(userService, authServiceClient)
	server := ioc.InitGrpcServer(userServiceServer)
	app := &App{
		server: server,
	}
	return app
}
