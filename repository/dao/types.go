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

type ArticleLikeDao interface {
	IncrReadCount(ctx context.Context, ObjType string, ObjID int64) error
	InSertLike(ctx context.Context, objType string, id int64, uid int64) error
	DeleteLike(ctx context.Context, objType string, id int64, uid int64) error
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

// TODO 找TOPn，超大规模数据情况下，归并

// Like 收集点赞数TOPn的数据
// 帖子的点赞，收藏，观看数
type Like struct {
	ID int64 `gorm:"primaryKey,autoIncrement"`

	// 联合索引， ObjID区分度更高，放左侧
	ObjID   int64  `gorm:"index:idx_objid_objtype"`
	ObjType string `gorm:"index:idx_objid_objtype;type:varchar(64)"`

	LikeCount    int64 `gorm:"column:like_count"`
	CollectCount int64 `gorm:"column:collect_count"`
	ViewCount    int64 `gorm:"column:view_count"`

	Ctime int64
	Utime int64
}

// UserGiveLike 用户点赞记录
type UserGiveLike struct {
	ID int64 `gorm:"primaryKey,autoIncrement"`

	// 此处的频繁查询应该是：查用户对帖子是否点赞，即where uid=? and objid=? and objtype=?
	// 考虑用户查看自己点赞的帖子，所以放左侧
	Uid     int64  `gorm:"uniqueIndex:idx_uid_objid_objtype"`
	ObjID   int64  `gorm:"uniqueIndex:idx_uid_objid_objtype"`
	ObjType string `gorm:"uniqueIndex:idx_uid_objid_objtype;type:varchar(64)"`

	Ctime int64
	Utime int64

	Status uint8 // 0：点赞，1：取消点赞
}
