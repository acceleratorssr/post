package ioc

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"time"
)

func InitOTEL() func(ctx context.Context) {
	// 创建并设置资源
	res, err := newResource("demo", "v0.0.1")
	if err != nil {
		panic(err)
	}

	// 配置传播器
	prop := newPropagator()
	// 在客户端和服务端之间传递 tracing 的相关信息
	otel.SetTextMapPropagator(prop)

	// 初始化跟踪提供者（trace provider）
	// 这个 provider 就是用来在打点的时候构建 trace 的
	//tracer := otel.Tracer("example-tracer")
	//ctx, span := tracer.Start(context.Background(), "example-operation")
	//time.Sleep(100 * time.Millisecond) // 模拟一些操作
	//span.End()
	tp, err := newTraceProvider(res)

	if err != nil {
		panic(err)
	}

	//defer tp.Shutdown(context.Background())
	otel.SetTracerProvider(tp)

	return func(ctx context.Context) {
		tp.Shutdown(ctx)
	}
}

func newResource(serviceName, serviceVersion string) (*resource.Resource, error) {
	return resource.Merge(resource.Default(),
		resource.NewWithAttributes(semconv.SchemaURL,
			semconv.ServiceName(serviceName),
			semconv.ServiceVersion(serviceVersion),
		))
}

// 可换成jaeger
func newTraceProvider(res *resource.Resource) (*trace.TracerProvider, error) {
	exporter, err := zipkin.New(
		"http://localhost:9411/api/v2/spans")
	if err != nil {
		return nil, err
	}

	traceProvider := trace.NewTracerProvider(
		trace.WithBatcher(exporter,
			trace.WithMaxExportBatchSize(20),
			// Default is 5s. Set to 1s for demonstrative purposes.
			trace.WithBatchTimeout(time.Second)),
		trace.WithResource(res),
	)
	return traceProvider, nil
}

func newPropagator() propagation.TextMapPropagator {
	return propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
}
