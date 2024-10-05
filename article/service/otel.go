package service

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"post/article/domain"
)

type OTELService struct {
	svc    articleService
	tracer trace.Tracer
}

func (O *OTELService) Publish(ctx context.Context, art *domain.Article) error {
	//TODO implement me
	panic("implement me")
}

func (O *OTELService) Withdraw(ctx context.Context, aid, uid uint64) error {
	//TODO implement me
	panic("implement me")
}

func (O *OTELService) ListSelf(ctx context.Context, uid uint64, list *domain.List) ([]domain.Article, error) {
	//TODO implement me
	panic("implement me")
}

func (O *OTELService) GetAuthorModelsByID(ctx context.Context, aid, uid uint64) (*domain.Article, error) {
	//TODO implement me
	panic("implement me")
}

func (O *OTELService) GetPublishedByID(ctx context.Context, id, uid uint64) (*domain.Article, error) {
	//TODO implement me
	panic("implement me")
}

func (O *OTELService) ListPublished(ctx context.Context, list *domain.List) ([]domain.Article, error) {
	//TODO implement me
	panic("implement me")
}

// Save 打点
func (O *OTELService) Save(ctx context.Context, art *domain.Article) error {
	ctx, span := O.tracer.Start(ctx, "articleService_save",
		trace.WithSpanKind(trace.SpanKindServer), // ...
	)
	defer span.End()
	span.AddEvent("save article")
	//// 从上下文中提取 Span
	//span := trace.SpanFromContext(ctx)

	//// grpc客户端传递，将 SpanContext 注入到 gRPC metadata 中
	//md := metadata.New(nil)
	//trace.Inject(ctx, otel.GetTextMapPropagator(), metadata.NewWriter(md))
	//ctx = metadata.NewOutgoingContext(ctx, md)

	// grpc服务端提取
	//if md, ok := metadata.FromIncomingContext(ctx); ok {
	//            ctx = trace.Extract(ctx, otel.GetTextMapPropagator(), metadata.NewReader(md))
	//        }
	//
	//        // 提取当前的 Span
	//        span := trace.SpanFromContext(ctx)
	//        defer span.End()

	err := O.svc.Save(ctx, art)
	if err != nil {
		span.RecordError(err)
	}
	return nil
}

// NewArticleServiceWithTracer 装饰器
func NewArticleServiceWithTracer(svc articleService) ArticleService {
	tp := otel.GetTracerProvider()
	tracer := tp.Tracer("internal/service/article.go")
	return &OTELService{
		tracer: tracer,
	}
}
