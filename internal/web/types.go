package web

import (
	"github.com/gin-gonic/gin"
	"post/internal/domain"
)

type handler interface {
	RegisterRoutes(s *gin.Engine)
}

type Req struct {
	ID      int64  `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

type ReqOnlyWithID struct {
	ID int64 `json:"id"`
}

type ReqList struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

type LikeReq struct {
	ObjID int64 `json:"obj_id"`
	Liked bool  `json:"liked"`
}

type CollectReq struct {
	Uid     int64  `json:"uid"`
	ObjID   int64  `json:"obj_id"`
	ObjType string `json:"obj_type"`
}

type ArticleVO struct {
	ID         int64  `json:"id"`
	Title      string `json:"title"`
	Content    string `json:"content"`
	AuthorID   int64  `json:"author_id"`
	AuthorName string `json:"author_name"`
	Status     uint8  `json:"status"`
	Ctime      string `json:"ctime"`
	Utime      string `json:"utime"`
	// 可考虑加上该三个字段，放缓存内，查文章时一起返回
	//ReadCount    int64  `json:"read_count"`
	//LikeCount    int64  `json:"like_count"`
	//CollectCount int64  `json:"collect_count"`
}

// Req 转换为 domain.Article
func (r Req) toDomain(id int64, name string) domain.Article {
	return domain.Article{
		ID:      r.ID,
		Title:   r.Title,
		Content: r.Content,
		Author: domain.Author{
			Id:   id,
			Name: name,
		},
	}
}

func (r ReqOnlyWithID) toDomain(id int64, name string) domain.Article {
	return domain.Article{
		ID: r.ID,
		Author: domain.Author{
			Id:   id,
			Name: name,
		},
	}
}
