//go:build wireinject

package main

import (
	"github.com/google/wire"
	"post/ioc"
	"post/repository"
	"post/repository/dao"
	"post/service"
	"post/web"
)

func InitApp() *web.ArticleHandler {
	wire.Build(
		//ioc.InitDB,
		ioc.InitMongoDB,
		//dao.NewGORMArticleDao,
		dao.NewMongoDB,
		dao.NewSnowflakeNode,
		repository.NewArticleAuthorRepository,
		repository.NewArticleReaderRepository,
		service.NewArticleService,
		web.NewArticleHandler,
		//ioc.InitWebServer,
	)
	return new(web.ArticleHandler)
}
