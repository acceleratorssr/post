//go:build wireinject

package main

import (
	"github.com/google/wire"
	"post/bff/ioc"
	"post/bff/web"
)

func InitApp() *App {
	wire.Build(
		ioc.InitLogger,
		ioc.InitEtcdClient,
		ioc.InitUserClient,
		ioc.InitSSOClient,
		ioc.NewJWTHandler,
		ioc.InitGinServer,

		web.NewUserHandler,
		web.NewSSOHandler,

		wire.Struct(new(App), "*"),
	)
	return new(App)
}
