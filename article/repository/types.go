package repository

import (
	"post/article/domain"
	"post/article/repository/dao"
)

func ToAuthorEntity(art *domain.Article) *dao.ArticleAuthor {
	return &dao.ArticleAuthor{
		SnowID:   art.ID,
		Title:    art.Title,
		Content:  art.Content,
		Authorid: art.Author.Id,
	}
}

func ToReaderEntity(art *domain.Article) *dao.ArticleReader {
	return &dao.ArticleReader{
		SnowID:   int64(art.ID),
		Title:    art.Title,
		Content:  art.Content,
		Authorid: art.Author.Id,
	}
}
