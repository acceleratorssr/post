package dao

import (
	"context"
	"post/migrator"
)

type ArticleLikeDao interface {
	IncrReadCount(ctx context.Context, ObjType string, ObjID int64) error
	IncrReadCountMany(ctx context.Context, ObjType string, ObjIDs []int64) error

	InSertLike(ctx context.Context, objType string, id int64, uid int64) error
	DeleteLike(ctx context.Context, objType string, id int64, uid int64) error
	InsertCollection(ctx context.Context, ObjType string, ObjID, uid int64) error

	GetPublishedByBatch(ctx context.Context, ObjType string, offset, limit int, now int64) ([]Like, error)
}

// TODO 找TOPn，超大规模数据情况下，归并

// Like 收集点赞数TOPn的数据
// 帖子的点赞，收藏，观看数
type Like struct {
	ID int64 `gorm_ex:"primaryKey,autoIncrement"`

	// 联合索引， ObjID区分度更高，放左侧
	ObjID   int64  `gorm_ex:"uniqueIndex:idx_objid_objtype"`
	ObjType string `gorm_ex:"uniqueIndex:idx_objid_objtype;type:varchar(64)"`

	LikeCount    int64 `gorm_ex:"column:like_count"`
	CollectCount int64 `gorm_ex:"column:collect_count"`
	ViewCount    int64 `gorm_ex:"column:view_count"`

	Ctime int64 `gorm_ex:"index:idx_ctime"`
	Utime int64
}

func (l Like) CompareWith(entity migrator.Entity) bool {
	e, ok := entity.(Like)
	if !ok {
		panic("entity is not \"Like\"")
	}
	return l == e
}

func (l Like) GetID() int64 {
	return l.ID
}

// UserGiveLike 用户点赞记录
type UserGiveLike struct {
	ID int64 `gorm_ex:"primaryKey,autoIncrement"`

	// 此处的频繁查询应该是：查用户对帖子是否点赞，即where uid=? and objid=? and objtype=?
	// 考虑用户查看自己点赞的帖子，所以放左侧
	Uid     int64  `gorm_ex:"uniqueIndex:idx_uid_objid_objtype"`
	ObjID   int64  `gorm_ex:"uniqueIndex:idx_uid_objid_objtype"`
	ObjType string `gorm_ex:"uniqueIndex:idx_uid_objid_objtype;type:varchar(64)"`

	Ctime int64
	Utime int64

	Status uint8 // 0：点赞，1：取消点赞
}

type UserGiveCollect struct {
	ID int64 `gorm_ex:"primaryKey,autoIncrement"`

	// 此处的频繁查询应该是：查用户对帖子是否点赞，即where uid=? and objid=? and objtype=?
	// 考虑用户查看自己点赞的帖子，所以放左侧
	Uid     int64  `gorm_ex:"uniqueIndex:idx_uid_objid_objtype"`
	ObjID   int64  `gorm_ex:"uniqueIndex:idx_uid_objid_objtype"`
	ObjType string `gorm_ex:"uniqueIndex:idx_uid_objid_objtype;type:varchar(64)"`

	Ctime int64
	Utime int64

	Status uint8 // 0：收藏，1：取消收藏
}

type UserGiveRead struct {
	ID int64 `gorm_ex:"primaryKey,autoIncrement"`

	Uid     int64  `gorm_ex:"uniqueIndex:idx_uid_objid_objtype"`
	ObjID   int64  `gorm_ex:"uniqueIndex:idx_uid_objid_objtype"`
	ObjType string `gorm_ex:"uniqueIndex:idx_uid_objid_objtype;type:varchar(64)"`

	Ctime int64
	Utime int64
}
