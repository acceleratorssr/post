package dao

import (
	"context"
)

type ArticleDAO interface {
	InputArticle(ctx context.Context, article Article) error
	Search(ctx context.Context, tagArtIds []int64, keywords []string) ([]Article, error)
}

type AnyDAO interface {
	Input(ctx context.Context, index, docID, data string) error
}
