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

func (s *OTELService) GetPublishedByIDS(ctx context.Context, ids []uint64) ([]domain.Article, error) {
	//TODO implement me
	panic("implement me")
}

// Save 打点
func (s *OTELService) Save(ctx context.Context, art *domain.Article) (uint64, error) {
	ctx, span := s.tracer.Start(ctx, "articleService_save",
		trace.WithSpanKind(trace.SpanKindServer), // ...
	)
	defer span.End()
	span.AddEvent("save article")

	id, err := s.svc.Save(ctx, art)
	if err != nil {
		span.RecordError(err)
	}
	return id, nil
}

func (s *OTELService) Publish(ctx context.Context, art *domain.Article) (uint64, error) {
	//TODO implement me
	panic("implement me")
}

func (s *OTELService) Withdraw(ctx context.Context, art *domain.Article) error {
	//TODO implement me
	panic("implement me")
}

func (s *OTELService) List(ctx context.Context, uid uint64, limit, offset int) ([]domain.Article, error) {
	//TODO implement me
	panic("implement me")
}

func (s *OTELService) GetAuthorModelsByID(ctx context.Context, id uint64) (*domain.Article, error) {
	//TODO implement me
	panic("implement me")
}

func (s *OTELService) GetPublishedByID(ctx context.Context, id, uid uint64) (*domain.Article, error) {
	//TODO implement me
	panic("implement me")
}

// NewArticleServiceWithTracer 装饰器
func NewArticleServiceWithTracer(svc articleService) *OTELService {
	tp := otel.GetTracerProvider()
	tracer := tp.Tracer("internal/service/article.go")
	return &OTELService{
		tracer: tracer,
	}
}
