package service

import (
	"context"
	"github.com/volcengine/volcengine-go-sdk/service/arkruntime"
	"github.com/volcengine/volcengine-go-sdk/service/arkruntime/model"
	"golang.org/x/sync/errgroup"
	"os"
	"post/search/domain"
	"post/search/repository"
	"strings"
)

type SearchService interface {
	Search(ctx context.Context, expression string, limit int) (domain.SearchResult, error)
}

type searchService struct {
	articleRepo repository.ArticleRepository
	marsCode    *arkruntime.Client
}

func NewSearchService(articleRepo repository.ArticleRepository, marsCode *arkruntime.Client) SearchService {
	return &searchService{
		articleRepo: articleRepo,
		marsCode:    marsCode,
	}
}

func (s *searchService) Search(ctx context.Context, expression string, limit int) (domain.SearchResult, error) {
	keywords := strings.FieldsFunc(expression, func(r rune) bool {
		// 判断哪些符号可以作为分隔符，支持空格、逗号、句号、分号等
		return r == ' ' || r == ',' || r == '.' || r == ';' || r == '，'
	})

	var embeddings []float32
	embeddingsResp, err := s.marsCode.CreateEmbeddings(ctx, model.EmbeddingRequestStrings{
		Input: keywords,
		Model: os.Getenv("EP_VECTOR_KEY"),
	})
	if err == nil && embeddingsResp.Data != nil {
		embeddings = embeddingsResp.Data[0].Embedding
	} else {
		// log
	}

	var eg errgroup.Group
	var res domain.SearchResult
	eg.Go(func() error {
		arts, err := s.articleRepo.SearchArticle(ctx, keywords, embeddings, limit)
		res.Articles = arts
		return err
	})
	return res, eg.Wait()
}
