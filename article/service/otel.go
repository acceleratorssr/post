package service

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"post/article/domain"
)

type ToOTelTracer ArticleService

type OTELService struct {
	svc    ToOTelTracer
	tracer trace.Tracer
}

func (O *OTELService) ListPublished(ctx context.Context, list *domain.List, uid uint64) ([]domain.Article, error) {
	ctx, span := O.tracer.Start(ctx, "article-service-listPublish",
		trace.WithSpanKind(trace.SpanKindServer), // ...
	)
	defer span.End()
	span.AddEvent("listPublish article")

	arts, err := O.svc.ListPublished(ctx, list, uid)
	if err != nil {
		span.RecordError(err)
	}
	return arts, nil
}

func (O *OTELService) GetArtByIDs(ctx context.Context, aids []uint64) ([]domain.Article, error) {
	ctx, span := O.tracer.Start(ctx, "article-service-getArtByIDs",
		trace.WithSpanKind(trace.SpanKindServer), // ...
	)
	defer span.End()
	span.AddEvent("getArtByIDs article")

	arts, err := O.svc.GetArtByIDs(ctx, aids)
	if err != nil {
		span.RecordError(err)
	}
	return arts, nil
}

func (O *OTELService) Publish(ctx context.Context, art *domain.Article) error {
	ctx, span := O.tracer.Start(ctx, "article-service-publish",
		trace.WithSpanKind(trace.SpanKindServer), // ...
	)
	defer span.End()
	span.AddEvent("publish article")

	err := O.svc.Publish(ctx, art)
	if err != nil {
		span.RecordError(err)
	}
	return nil
}

func (O *OTELService) Withdraw(ctx context.Context, aid, uid uint64) error {
	ctx, span := O.tracer.Start(ctx, "article-service-withdraw",
		trace.WithSpanKind(trace.SpanKindServer), // ...
	)
	defer span.End()
	span.AddEvent("withdraw article")

	err := O.svc.Withdraw(ctx, aid, uid)
	if err != nil {
		span.RecordError(err)
	}
	return nil
}

func (O *OTELService) ListSelf(ctx context.Context, uid uint64, list *domain.List) ([]domain.Article, error) {
	ctx, span := O.tracer.Start(ctx, "article-service-listSelf",
		trace.WithSpanKind(trace.SpanKindServer), // ...
	)
	defer span.End()
	span.AddEvent("listSelf article")

	arts, err := O.svc.ListSelf(ctx, uid, list)
	if err != nil {
		span.RecordError(err)
	}
	return arts, nil
}

func (O *OTELService) GetAuthorModelsByID(ctx context.Context, aid, uid uint64) (*domain.Article, error) {
	ctx, span := O.tracer.Start(ctx, "article-service-getAuthorModelByID",
		trace.WithSpanKind(trace.SpanKindServer), // ...
	)
	defer span.End()
	span.AddEvent("getAuthorModelByID article")

	art, err := O.svc.GetAuthorModelsByID(ctx, aid, uid)
	if err != nil {
		span.RecordError(err)
	}
	return art, nil
}

func (O *OTELService) GetPublishedByID(ctx context.Context, id, uid uint64) (*domain.Article, error) {
	ctx, span := O.tracer.Start(ctx, "article-service-getPublishedByID",
		trace.WithSpanKind(trace.SpanKindServer), // ...
	)
	defer span.End()
	span.AddEvent("getPublishedByID article")

	art, err := O.svc.GetPublishedByID(ctx, id, uid)
	if err != nil {
		span.RecordError(err)
	}
	return art, nil
}

// Save 打点
func (O *OTELService) Save(ctx context.Context, art *domain.Article) error {
	ctx, span := O.tracer.Start(ctx, "article-service-save",
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

func NewArticleServiceWithTracer(svc ToOTelTracer) ArticleService {
	tp := otel.GetTracerProvider()
	tracer := tp.Tracer("service/articles")
	return &OTELService{
		svc:    svc,
		tracer: tracer,
	}
}
