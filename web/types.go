package web

import (
	"github.com/gin-gonic/gin"
	"post/domain"
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
	ID    int64 `json:"obj_id"`
	Liked bool  `json:"liked"`
}

type ArticleVO struct {
	ID         int64
	Title      string
	Content    string
	AuthorID   int64
	AuthorName string
	Status     uint8
	Ctime      string
	Utime      string
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
