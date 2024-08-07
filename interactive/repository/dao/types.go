package dao

import "context"

type ArticleLikeDao interface {
	IncrReadCount(ctx context.Context, ObjType string, ObjID int64) error
	IncrReadCountMany(ctx context.Context, ObjType string, ObjIDs []int64) error

	InSertLike(ctx context.Context, objType string, id int64, uid int64) error
	DeleteLike(ctx context.Context, objType string, id int64, uid int64) error
	InsertCollection(ctx context.Context, ObjType string, ObjID, uid int64) error

	GetPublishedByBatch(ctx context.Context, ObjType string, offset, limit int, now int64) ([]Like, error)
}
