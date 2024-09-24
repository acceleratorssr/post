package integration

import (
	"context"
	"encoding/json"
	"github.com/olivere/elastic/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"log"
	searchv1 "post/api/proto/gen/search/v1"
	"post/search/grpc"
	"post/search/integration/startup"
	"testing"
	"time"
)

type SearchTestSuite struct {
	suite.Suite
	searchSvc *grpc.SearchServiceServer
	syncSvc   *grpc.SyncServiceServer
}

func (s *SearchTestSuite) SetupSuite() {
	s.searchSvc = startup.InitSearchServer()
	s.syncSvc = startup.InitSyncServer()
}

func (s *SearchTestSuite) TestSearch() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	data, err := json.Marshal(Tags{
		ObjType: "article",
		ObjId:   123,
		Tags:    []string{"yes"},
	})
	require.NoError(s.T(), err)
	_, err = s.syncSvc.InputAny(ctx, &searchv1.InputAnyRequest{
		IndexName: "tags_index",
		DocId:     "abcd",
		Data:      string(data),
	})
	require.NoError(s.T(), err)

	_, err = s.syncSvc.InputArticle(ctx, &searchv1.InputArticleRequest{
		Article: &searchv1.Article{
			Id:     123,
			Title:  "title2,yy",
			Status: 2,
		},
	})
	require.NoError(s.T(), err)
	_, err = s.syncSvc.InputArticle(ctx, &searchv1.InputArticleRequest{
		Article: &searchv1.Article{
			Id:      222,
			Content: "content3,歪比八不",
			Status:  2,
		},
	})
	require.NoError(s.T(), err)
	resp, err := s.searchSvc.Search(ctx, &searchv1.SearchRequest{
		Expression: "八不 yes",
	})
	require.NoError(s.T(), err)
	assert.Equal(s.T(), 2, len(resp.Article.Articles))
	//over()
}

type Tags struct {
	ObjType string   `json:"obj_type"`
	ObjId   int64    `json:"obj_id"`
	Tags    []string `json:"tags"`
}

func TestSearchService(t *testing.T) {
	suite.Run(t, new(SearchTestSuite))
}

func over() {
	// 初始化 Elasticsearch 客户端
	const timeout = 10 * time.Second
	opts := []elastic.ClientOptionFunc{
		elastic.SetURL("http://localhost:9200"),
		elastic.SetSniff(false),
		elastic.SetHealthcheckTimeoutStartup(timeout),
		elastic.SetTraceLog(log.Default()),
	}
	client, _ := elastic.NewClient(opts...)

	// 执行删除所有索引的请求
	_, _ = client.DeleteIndex("_all").Do(context.Background())
}
