package es_extra

import (
	"context"
	"fmt"
	es "github.com/elastic/go-elasticsearch/v8"
	"testing"
)

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
					ID   int64  `json:"id"`
					Name string `json:"name"`
				} `json:"author"`
				Embedding []float32 `json:"embedding"`
			} `json:"_source"`
			Score float64 `json:"_score"`
		} `json:"hits"`
	} `json:"hits"`
}

func TestQuery(t *testing.T) {
	esClient, _ := es.NewDefaultClient()
	ctx := context.Background()

	builder := NewQueryBuilder().
		MultiMatch("backend", []string{"title", "content"})

	search, err := Search[SearchResponse](ctx, esClient, builder, 10)
	if err != nil || search.Hits.Hits == nil {
		t.Errorf("Search failed: %v", err)
	}
	//fmt.Println(search)

	builder = NewQueryBuilder().MultiMatch("backend", []string{"title", "content"})
	//b2 := NewQueryBuilder().Fuzzy("title", "haalth", "AUTO", 5.0)
	b2 := NewQueryBuilder().CosineSimilarity("embedding", []float32{0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8, 0.9, 1.0})
	boolBuilder := NewBoolQuery().Should(builder, b2)

	search, err = Search[SearchResponse](ctx, esClient, boolBuilder, 10)
	if err != nil || search.Hits.Hits == nil {
		t.Errorf("Search failed: %v", err)
	}
	fmt.Println(search)
}
