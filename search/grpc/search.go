package grpc

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	searchv1 "post/api/proto/gen/search/v1"
	"post/search/domain"
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
		return nil, status.Errorf(codes.Internal, "interactive 搜索文章失败: %s", err)
	}

	return &searchv1.SearchResponse{
		Article: &searchv1.ArticleResult{
			Articles: s.toDTO(resp.Articles...),
		},
	}, nil
}

func (s *SearchServiceServer) toDTO(art ...domain.Article) []*searchv1.Article {
	articles := make([]*searchv1.Article, len(art))

	for i, src := range art {
		articles[i] = &searchv1.Article{
			Id:      src.ID,
			Title:   src.Title,
			Content: src.Content,
		}
	}
	return articles
}
