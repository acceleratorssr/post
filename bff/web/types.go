package web

import (
	"github.com/gin-gonic/gin"
	"post/article/domain"
)

type handler interface {
	RegisterRoutes(engine *gin.Engine, mw gin.HandlerFunc)
}

type Req struct {
	ID      uint64 `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

type ReqOnlyWithID struct {
	ID uint64 `json:"id"`
}

type ReqList struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

type LikeReq struct {
	ObjID uint64 `json:"obj_id"`
	Liked bool   `json:"liked"`
}

type CollectReq struct {
	Uid     uint64 `json:"uid"`
	ObjID   uint64 `json:"obj_id"`
	ObjType string `json:"obj_type"`
}

type ArticleResp struct {
	ID         uint64 `json:"id"`
	Title      string `json:"title"`
	Content    string `json:"content"`
	AuthorID   uint64 `json:"author_id"`
	AuthorName string `json:"author_name"`
	Status     uint8  `json:"status"`
	Ctime      int64  `json:"ctime"`
	Utime      int64  `json:"utime"`
	// 可考虑加上该三个字段，放缓存内，查文章时一起返回
	//ReadCount    int64  `json:"read_count"`
	//LikeCount    int64  `json:"like_count"`
	//CollectCount int64  `json:"collect_count"`
}

// Req 转换为 domain.Article
func (r Req) toDomain(id uint64, name string) *domain.Article {
	return &domain.Article{
		ID:      r.ID,
		Title:   r.Title,
		Content: r.Content,
		Author: domain.Author{
			Id:   id,
			Name: name,
		},
	}
}

func (r ReqOnlyWithID) toDomain(id uint64, name string) *domain.Article {
	return &domain.Article{
		ID: r.ID,
		Author: domain.Author{
			Id:   id,
			Name: name,
		},
	}
}
