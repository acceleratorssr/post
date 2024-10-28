package dao

import (
	"context"
)

type ArticleDAO interface {
	InputArticle(ctx context.Context, article Article, vector []float32) error
	Search(ctx context.Context, tagArtIds []int64, keywords []string, vector []float32, limit int) ([]Article, error)
	DeleteArticle(ctx context.Context, id uint64) error
}

type AnyDAO interface {
	Input(ctx context.Context, index, docID, data string) error
}

type TagDAO interface {
	Search(ctx context.Context, objType string, keywords []string) ([]int64, error)
}
