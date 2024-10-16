//go:build wireinject

package main

import (
	"github.com/google/wire"
	"post/recommend/events"
	"post/recommend/grpc"
	"post/recommend/ioc"
	"post/recommend/service"
)

func InitApp() *App {
	wire.Build(
		ioc.InitKafka,
		ioc.InitGorse,

		events.NewKafkaRecommendConsumer,

		service.NewRecommendService,
		grpc.NewRecommendServiceServer,

		ioc.InitGrpcRecommendServer,
		ioc.NewConsumers,
		wire.Struct(new(App), "*"),
	)
	return new(App)
}
