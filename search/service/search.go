package service

import (
	"context"
	"golang.org/x/sync/errgroup"
	"post/search/domain"
	"post/search/repository"
	"strings"
)

type SearchService interface {
	Search(ctx context.Context, expression string) (domain.SearchResult, error)
}

type searchService struct {
	articleRepo repository.ArticleRepository
}

func NewSearchService(articleRepo repository.ArticleRepository) SearchService {
	return &searchService{articleRepo: articleRepo}
}

func (s *searchService) Search(ctx context.Context, expression string) (domain.SearchResult, error) {
	keywords := strings.FieldsFunc(expression, func(r rune) bool {
		// 判断哪些符号可以作为分隔符，支持空格、逗号、句号、分号等
		return r == ' ' || r == ',' || r == '.' || r == ';' || r == '，'
	})

	var eg errgroup.Group
	var res domain.SearchResult
	eg.Go(func() error {
		arts, err := s.articleRepo.SearchArticle(ctx, keywords)
		res.Articles = arts
		return err
	})
	return res, eg.Wait()
}
