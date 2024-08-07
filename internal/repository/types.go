package repository

import (
	"post/internal/domain"
	"post/internal/repository/dao"
)

func ToEntity(art domain.Article) dao.ArticleAuthor {
	return dao.ArticleAuthor{
		Id:       art.ID,
		Title:    art.Title,
		Content:  art.Content,
		Authorid: art.Author.Id,
		Status:   art.Status.ToUint8(),
	}
}
