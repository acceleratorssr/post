// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"post/bff/ioc"
	"post/bff/web"
)

// Injectors from wire.go:

func InitApp() *App {
	logger := ioc.InitLogger()
	client := ioc.InitEtcdClient()
	authServiceClient := ioc.InitSSOClient(client)
	jwt := ioc.NewJWTHandler(authServiceClient)
	userServiceClient := ioc.InitUserClient(client)
	userHandler := web.NewUserHandler(userServiceClient, authServiceClient)
	ssoHandler := web.NewSSOHandler(authServiceClient)
	articleServiceClient := ioc.InitArticleClient(client)
	likeServiceClient := ioc.InitInteractiveClient(client)
	articleHandler := web.NewArticleHandler(articleServiceClient, likeServiceClient)
	server := ioc.InitGinServer(logger, jwt, userHandler, ssoHandler, articleHandler)
	app := &App{
		WebServer: server,
	}
	return app
}