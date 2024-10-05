package dao

import (
	"context"
)

type ArticleDao interface {
	Insert(ctx context.Context, art *ArticleAuthor) error
	UpsertReader(ctx context.Context, art *ArticleReader) error

	UpdateByID(ctx context.Context, art *ArticleAuthor) error
	DeleteReader(ctx context.Context, aid uint64, uid uint64) error
	GetListByAuthor(ctx context.Context, uid uint64, list *List) ([]ArticleAuthor, error)
	GetAuthorByID(ctx context.Context, aid, uid uint64) (*ArticleAuthor, error)
	GetPublishedByID(ctx context.Context, id uint64) (*ArticleReader, error)
	ListPublished(ctx context.Context, list *List) ([]ArticleReader, error)
	ListByID(ctx context.Context, uid uint64, list *List) ([]ArticleReader, error)
}

// ArticleAuthor 为ing库，或者说草稿库
type ArticleAuthor struct {
	SnowID     uint64 `gorm:"primaryKey"`
	Title      string `gorm:"size:2048"`
	Content    string `gorm:"type=BLOB"`
	AuthorName string `gorm:"type:varchar(64);size:64"`
	Authorid   uint64 `gorm:"index"`
	Ctime      int64  `gorm:"index"`
	Utime      int64
}

type ArticleReader struct {
	ID       uint64 `gorm:"primaryKey,autoIncrement"`
	Title    string `gorm:"size:2048"`
	Content  string `gorm:"type=BLOB"`
	Authorid uint64 `gorm:"index"`
	Ctime    int64  `gorm:"index"`
	Utime    int64

	SnowID int64 `gorm:"uniqueIndex"`
}

// 使用callback代替hook
// https://gorm.io/zh_CN/docs/write_plugins.html
//func (aa *ArticleAuthor) BeforeCreate(tx *gorm_ex.DB) (err error) {
//	start := time.Now()
//	tx.Set("start", start)
//	return
//}
//
//func (aa *ArticleAuthor) AfterCreate(tx *gorm_ex.DB) (err error) {
//	start, ok := tx.Statement.Get("start")
//	if !ok {
//		return
//	}
//	t, ok := start.(time.Time)
//	if !ok {
//		return
//	}
//	fmt.Println(time.Since(t))
//	return
//}

type List struct {
	Limit     int
	LastValue int64 // 保存在客户端，用于翻页时防重复数据
	Desc      bool  // 0为升序，1为降序
	OrderBy   string
}
