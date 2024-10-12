package web

import (
	"github.com/gin-gonic/gin"
	searchv1 "post/api/proto/gen/search/v1"
	gin_extra "post/pkg/gin-extra"
	"post/search/domain"
)

type SearchHandler struct {
	svc searchv1.SearchServiceClient
}

type SearchReq struct {
	Expression string `json:"expression"`
}

func (s *SearchHandler) RegisterRoutes(engine *gin.Engine, mw gin.HandlerFunc) {
	search := engine.Group("/search")
	search.Use(mw)
	search.POST("/article", gin_extra.WrapWithReq[SearchReq](s.Search))
}

func (s *SearchHandler) Search(ctx *gin.Context, req SearchReq) (*gin_extra.Response, error) {
	uid, ok := ctx.Get("uid")
	if !ok {
		return &gin_extra.Response{
			Code: gin_extra.Unauthenticated,
			Msg:  "请先登录",
		}, nil
	}
	search, err := s.svc.Search(ctx, &searchv1.SearchRequest{
		Uid:        uid.(uint64),
		Expression: req.Expression,
	})
	if err != nil {
		return &gin_extra.Response{
			Code: gin_extra.System,
		}, err
	}

	return &gin_extra.Response{
		Data: s.toDTO(search.Article.Articles...),
	}, nil
}

func (s *SearchHandler) toDTO(art ...*searchv1.Article) []*domain.Article {
	articles := make([]*domain.Article, len(art))
	for i, src := range art {
		articles[i] = &domain.Article{
			ID:      src.Id,
			Title:   src.Title,
			Content: src.Content,
		}
	}
	return articles
}

func NewSearchHandler(svc searchv1.SearchServiceClient) *SearchHandler {
	return &SearchHandler{
		svc: svc,
	}
}
