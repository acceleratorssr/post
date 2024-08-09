package repository

import (
	"post/internal/domain"
	"post/internal/repository/dao"
)

func ToAuthorEntity(art domain.Article) dao.ArticleAuthor {
	return dao.ArticleAuthor{
		Id:       art.ID,
		Title:    art.Title,
		Content:  art.Content,
		Authorid: art.Author.Id,
		Status:   art.Status.ToUint8(),
	}
}

func ToReaderEntity(art domain.Article) dao.ArticleReader {
	return dao.ArticleReader{
		Id:       art.ID,
		Title:    art.Title,
		Content:  art.Content,
		Authorid: art.Author.Id,
		Status:   art.Status.ToUint8(),
	}
}
