//go:build wireinject

package main

import (
	"github.com/google/wire"
	"post/ioc"
	"post/repository"
	"post/repository/cache"
	"post/repository/dao"
	"post/service"
	"post/web"
)

func InitApp() *web.ArticleHandler {
	wire.Build(
		//ioc.InitDB,
		ioc.InitMongoDB,
		ioc.InitLogger,

		//dao.NewGORMArticleDao,
		dao.NewMongoDB,
		cache.NewRedisArticleCache,
		dao.NewSnowflakeNode,

		repository.NewArticleAuthorRepository,
		repository.NewArticleReaderRepository,

		service.NewArticleService,
		service.NewLikeService,

		web.NewArticleHandler,
		//ioc.InitWebServer,
	)
	return new(web.ArticleHandler)
}
