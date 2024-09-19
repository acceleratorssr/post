package service

import (
	"context"
	"post/search/domain"
	"post/search/repository"
)

type SyncService interface {
	InputArticle(ctx context.Context, article domain.Article) error
	InputAny(ctx context.Context, index, docID, data string) error
}

type syncService struct {
	articleRepo repository.ArticleRepository
	anyRepo     repository.AnyRepository
}

func (s *syncService) InputAny(ctx context.Context, index, docID, data string) error {
	return s.anyRepo.Input(ctx, index, docID, data)
}

func (s *syncService) InputArticle(ctx context.Context, article domain.Article) error {
	return s.articleRepo.InputArticle(ctx, article)
}

func NewSyncService(
	anyRepo repository.AnyRepository,
	articleRepo repository.ArticleRepository) SyncService {
	return &syncService{
		articleRepo: articleRepo,
		anyRepo:     anyRepo,
	}
}
