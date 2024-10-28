package dao

import (
	"context"
	"encoding/json"
	"fmt"
	es "github.com/elastic/go-elasticsearch/v8"
	"post/pkg/es-extra"
	esx "post/pkg/es-extra"
	"strconv"
	"strings"
)

type Article struct {
	Id        uint64    `json:"id"`
	Title     string    `json:"title"`
	Status    int32     `json:"status"`
	Content   string    `json:"content"`
	Tags      []string  `json:"tags"`
	Author    Author    `json:"author"`
	Embedding []float32 `json:"embedding"`
}

type Author struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
}

type SearchResponse struct {
	Hits struct {
		Hits []struct {
			Source struct {
				ID          int64    `json:"id"`
				Title       string   `json:"title"`
				Status      int      `json:"status"`
				Content     string   `json:"content"`
				Tags        []string `json:"tags"`
				PublishedAt string   `json:"published_at"`
				Author      struct {
					ID   uint64 `json:"id"`
					Name string `json:"name"`
				} `json:"author"`
				Embedding []float32 `json:"embedding"`
			} `json:"_source"`
			Score float64 `json:"_score"`
		} `json:"hits"`
	} `json:"hits"`
}

type ArticleElasticDAO struct {
	client *es.Client
}

type Searcher[T any] struct {
	client  *es.Client
	idxName []string
}

func NewSearcher[T any](client *es.Client, idxName []string) *Searcher[T] {
	return &Searcher[T]{
		client:  client,
		idxName: idxName,
	}
}

func NewArticleElasticDAO(client *es.Client) ArticleDAO {
	return &ArticleElasticDAO{client: client}
}

// Search todo 优化可选参数等 搜索的优化
func (h *ArticleElasticDAO) Search(ctx context.Context, tagArtIds []int64, keywords []string, vector []float32, limit int) ([]Article, error) {
	queryString := strings.Join(keywords, " ")

	id := es_extra.NewQueryBuilder().MultiMatch(queryString, []string{"title", "content"})

	v := es_extra.NewQueryBuilder().CosineSimilarity("embedding", vector)
	//v := es_extra.NewQueryBuilder().DotProduct("embedding", vector)
	//v := es_extra.NewQueryBuilder().EuclideanDistance("embedding", vector) // 效果差
	//v := es_extra.NewQueryBuilder().ManhattanDistance("embedding", vector) // 效果差
	//v := es_extra.NewQueryBuilder().KNN("embedding", vector, limit, 2*limit)// 效果差

	builder := es_extra.NewBoolQuery().Should(id, v)

	search, err := es_extra.Search[SearchResponse](ctx, h.client, builder, limit)
	if err != nil {
		return nil, err
	}

	arts := make([]Article, len(search.Hits.Hits))
	// 可设定分数阈值
	for i, hit := range search.Hits.Hits {
		arts[i] = h.toArt(&hit)
		fmt.Println(hit.Score)
	}

	return arts, nil
}

// InputArticle upsert
func (h *ArticleElasticDAO) InputArticle(ctx context.Context, article Article, vector []float32) error {
	article.Embedding = vector
	doc, _ := json.Marshal(article)
	res, err := h.client.Index(esx.ArticleIndexName, strings.NewReader(string(doc)),
		h.client.Index.WithRefresh("true"),
		h.client.Index.WithDocumentID(strconv.FormatUint(article.Id, 10)),
		h.client.Index.WithContext(ctx))
	if err != nil {
		// log
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		fmt.Println(res)
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			// log
		} else {
			// log
			//fmt.Printf("failed to index document: [%s] %s: %s",
			//	res.Status(),
			//	e["error"].(map[string]interface{})["type"],
			//	e["error"].(map[string]interface{})["reason"],
			//)
		}
		return fmt.Errorf("indexing error: %s", res.Status())
	}

	return nil
}

func (h *ArticleElasticDAO) DeleteArticle(ctx context.Context, articleID uint64) error {
	res, err := h.client.Delete(esx.ArticleIndexName,
		strconv.FormatUint(articleID, 10),
		h.client.Delete.WithContext(ctx),
	)
	if err != nil {
		// log
		return err
	}

	defer res.Body.Close()
	if res.IsError() {
		// log
		return fmt.Errorf("deletion error: %s", res.Status())
	}

	return nil
}

func (h *ArticleElasticDAO) toArt(hit *struct {
	Source struct {
		ID          int64    `json:"id"`
		Title       string   `json:"title"`
		Status      int      `json:"status"`
		Content     string   `json:"content"`
		Tags        []string `json:"tags"`
		PublishedAt string   `json:"published_at"`
		Author      struct {
			ID   uint64 `json:"id"`
			Name string `json:"name"`
		} `json:"author"`
		Embedding []float32 `json:"embedding"`
	} `json:"_source"`
	Score float64 `json:"_score"`
}) Article {
	return Article{
		Id:      uint64(hit.Source.ID),
		Title:   hit.Source.Title,
		Status:  int32(hit.Source.Status),
		Content: hit.Source.Content,
		Tags:    hit.Source.Tags,
	}
}
