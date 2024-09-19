package service

import (
	"context"
	"golang.org/x/sync/errgroup"
	"post/search/domain"
	"post/search/repository"
	"strings"
)

type SearchService interface {
	Search(ctx context.Context, uid int64, expression string) (domain.SearchResult, error)
}

type searchService struct {
	articleRepo repository.ArticleRepository
}

func NewSearchService(articleRepo repository.ArticleRepository) SearchService {
	return &searchService{articleRepo: articleRepo}
}

func (s *searchService) Search(ctx context.Context, uid int64, expression string) (domain.SearchResult, error) {
	// 这边一般要对 expression 进行一些预处理
	// 正常大家都是使用的空格符来分割的，但是有些时候可能会手抖，输错
	keywords := strings.Split(expression, " ")
	// 注意这里我们没有使用 multi query 或者 multi match 之类的写法
	// 是因为正常来说，不同的业务放过来的数据，什么支持搜索，什么不支持搜索，
	// 以及究竟怎么用于搜索，都是有区别的。所以这里我们利用两个 repo 来组合结果
	var eg errgroup.Group
	var res domain.SearchResult
	eg.Go(func() error {
		arts, err := s.articleRepo.SearchArticle(ctx, uid, keywords)
		res.Articles = arts
		return err
	})
	return res, eg.Wait()
}
