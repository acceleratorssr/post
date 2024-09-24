package grpc

import (
	"context"
	"google.golang.org/grpc"
	searchv1 "post/api/proto/gen/search/v1"
	"post/search/service"
)

type SearchServiceServer struct {
	searchv1.UnimplementedSearchServiceServer
	svc service.SearchService
}

func NewSearchService(svc service.SearchService) *SearchServiceServer {
	return &SearchServiceServer{svc: svc}
}

func (s *SearchServiceServer) Register(server grpc.ServiceRegistrar) {
	searchv1.RegisterSearchServiceServer(server, s)
}

func (s *SearchServiceServer) Search(ctx context.Context, request *searchv1.SearchRequest) (*searchv1.SearchResponse, error) {
	resp, err := s.svc.Search(ctx, request.Expression)
	if err != nil {
		return nil, err
	}

	articles := make([]*searchv1.Article, len(resp.Articles))

	for i, src := range resp.Articles {
		articles[i] = &searchv1.Article{
			Id:      src.Id,
			Title:   src.Title,
			Status:  src.Status,
			Content: src.Content,
		}
	}

	return &searchv1.SearchResponse{
		Article: &searchv1.ArticleResult{
			Articles: articles,
		},
	}, nil
}
