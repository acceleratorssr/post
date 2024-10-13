package dao

import (
	"context"
	"encoding/json"
	"github.com/olivere/elastic/v7"
	"strconv"
	"strings"
)

const ArticleIndexName = "article_index"
const TagIndexName = "tags_index"

type Article struct {
	Id      uint64   `json:"id"`
	Title   string   `json:"title"`
	Status  int32    `json:"status"`
	Content string   `json:"content"`
	Tags    []string `json:"tags"`
}

type ArticleElasticDAO struct {
	client *elastic.Client
}

type Searcher[T any] struct {
	client  *elastic.Client
	idxName []string
	query   elastic.Query
}

func NewSearcher[T any](client *elastic.Client, idxName []string, query elastic.Query) *Searcher[T] {
	return &Searcher[T]{
		client:  client,
		idxName: idxName,
		query:   query,
	}
}

func NewArticleElasticDAO(client *elastic.Client) ArticleDAO {
	return &ArticleElasticDAO{client: client}
}

func (h *ArticleElasticDAO) Search(ctx context.Context, tagArtIds []int64, keywords []string) ([]Article, error) {
	queryString := strings.Join(keywords, " ")
	ids := make([]interface{}, len(tagArtIds))
	for idx, src := range tagArtIds {
		ids[idx] = src
	}

	id := elastic.NewTermsQuery("id", ids...).Boost(2)   // 加权重
	title := elastic.NewMatchQuery("title", queryString) // 模糊匹配
	content := elastic.NewMatchQuery("content", queryString)
	status := elastic.NewTermQuery("status", 2)

	or := elastic.NewBoolQuery().Should( // or
		//elastic.NewTermsQuery("id", ids...).Boost(2), // 精确匹配
		title, content,
	)
	if len(ids) > 0 { // 避免tags没查出数据导致id为空
		or.Should(id)
	}
	query := elastic.NewBoolQuery().Must(or, status)

	return NewSearcher[Article](h.client, []string{ArticleIndexName}, query).Query(query).Do(ctx)
}

func (s *Searcher[T]) Query(query elastic.Query) *Searcher[T] {
	s.query = query
	return s
}

func (s *Searcher[T]) Do(ctx context.Context) ([]T, error) {
	resp, err := s.client.Search(s.idxName...).Query(s.query).Do(ctx)
	if err != nil {
		return nil, err
	}

	res := make([]T, 0, len(resp.Hits.Hits))
	for _, hit := range resp.Hits.Hits {
		var ele T
		_ = json.Unmarshal(hit.Source, &ele)
		res = append(res, ele)
	}
	return res, nil
}

// InputArticle upsert
func (h *ArticleElasticDAO) InputArticle(ctx context.Context, art Article) error {
	_, err := h.client.Index().
		Index(ArticleIndexName).
		Id(strconv.FormatUint(art.Id, 10)).
		BodyJson(art).Do(ctx)

	return err
}

func (h *ArticleElasticDAO) DeleteArticle(ctx context.Context, articleID uint64) error {
	_, err := h.client.Delete().
		Index(ArticleIndexName).
		Id(strconv.FormatUint(articleID, 10)).
		Do(ctx)

	return err
}
