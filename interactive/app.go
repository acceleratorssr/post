package main

import (
	"post/pkg/gin_ex"
	"post/pkg/grpc_ex"
	"post/pkg/sarama_ex"
)

// App 控制main中的方法的启用，控制生命周期
type App struct {
	server    *grpc_ex.Server
	consumers []sarama_ex.Consumer
	webAdmin  *gin_ex.Server
}
