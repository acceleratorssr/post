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
		ioc.InitInteractiveClient,
		ioc.InitArticleClient,
		ioc.InitSearchClient,

		ioc.NewJWTHandler,
		ioc.InitGinServer,

		web.NewUserHandler,
		web.NewSSOHandler,
		web.NewArticleHandler,
		web.NewSearchHandler,

		wire.Struct(new(App), "*"),
	)
	return new(App)
}
