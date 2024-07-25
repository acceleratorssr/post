package service

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"post/domain"
)

type Service struct {
	svc    articleService
	tracer trace.Tracer
}

func NewService(svc articleService) *Service {
	tp := otel.GetTracerProvider()
	tracer := tp.Tracer("post/service/article.go")
	return &Service{
		tracer: tracer,
	}
}

// Save 打点
func (s Service) Save(ctx context.Context, art domain.Article) (int64, error) {
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

func (s Service) Publish(ctx context.Context, art domain.Article) (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (s Service) Withdraw(ctx context.Context, art domain.Article) error {
	//TODO implement me
	panic("implement me")
}

func (s Service) List(ctx context.Context, uid int64, limit, offset int) ([]domain.Article, error) {
	//TODO implement me
	panic("implement me")
}

func (s Service) GetAuthorModelsByID(ctx context.Context, id int64) (domain.Article, error) {
	//TODO implement me
	panic("implement me")
}

func (s Service) GetPublishedByID(ctx context.Context, id, uid int64) (domain.Article, error) {
	//TODO implement me
	panic("implement me")
}
