package dao

import (
	"context"
)

type ArticleDao interface {
	Insert(ctx context.Context, art ArticleAuthor) (int64, error)
	UpdateByID(ctx context.Context, art ArticleAuthor) error
	SyncStatus(ctx context.Context, art ArticleAuthor) error
	GetListByAuthor(ctx context.Context, uid int64, limit int, offset int) ([]ArticleAuthor, error)
	GetByID(ctx context.Context, id int64) (ArticleAuthor, error)
	GetPublishedByID(ctx context.Context, id int64) (ArticleReader, error)
}

// ArticleAuthor 为ing库，或者说草稿库，制作库
type ArticleAuthor struct {
	Id       int64  `gorm:"primaryKey,autoIncrement"`
	Title    string `gorm:"size:4096"`
	Content  string `gorm:"type=BLOB"`
	Authorid int64  `gorm:"index"`
	Ctime    int64  `gorm:"index"`
	Utime    int64
	Status   uint8
}

type ArticleReader struct {
	Id       int64  `gorm:"primaryKey,autoIncrement"`
	Title    string `gorm:"size:4096"`
	Content  string `gorm:"type=BLOB"`
	Authorid int64  `gorm:"index"`
	Ctime    int64  `gorm:"index"`
	Utime    int64
	Status   uint8
}
