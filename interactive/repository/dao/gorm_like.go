package dao

import (
	"context"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"time"
)

type GORMArticleLikeDao struct {
	db *gorm.DB
}

func NewGORMArticleLikeDao(db *gorm.DB) ArticleLikeDao {
	return &GORMArticleLikeDao{
		db: db,
	}
}

// IncrReadCountMany todo 考虑到可以提前合并ObjIDs，在计数时可以直接+n
func (gad *GORMArticleLikeDao) IncrReadCountMany(ctx context.Context, ObjType string, ObjIDs []int64) error {
	now := time.Now().UnixMilli()
	likes := make([]Like, 0, len(ObjIDs))
	for i := 0; i < len(ObjIDs); i++ {
		likes = append(likes, Like{
			ObjID:     ObjIDs[i],
			ObjType:   ObjType,
			ViewCount: 1,
			Ctime:     now,
			Utime:     now,
		})
	}

	////此处只有一个事务，即事务提交仅触发一次刷入磁盘，批量只消耗一次磁盘IO
	//return gad.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
	//	txDAO := NewGORMArticleLikeDao(tx)
	//	for i := 0; i < len(ObjIDs); i++ {
	//		err := txDAO.IncrReadCount(ctx, ObjType, ObjIDs[i])
	//		if err != nil {
	//			return err
	//		}
	//	}
	//	return nil
	//})

	//return gad.db.WithContext(ctx).Clauses(clause.OnConflict{
	//	Columns: []clause.Column{{Name: "obj_id"}, {Name: "obj_type"}}, // 指定冲突列
	//	DoUpdates: clause.Assignments(map[string]any{
	//		"view_count": gorm.Expr("view_count + ?", 1),
	//		"utime":      time.Now().UnixMilli(),
	//	}),
	//}).CreateInBatches(likes, len(ObjIDs)).Error // len(ObjIDs)是每批次插入的记录数
	tx := gad.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "obj_id"}, {Name: "obj_type"}}, // 指定冲突列
		DoUpdates: clause.Assignments(map[string]any{
			"view_count": gorm.Expr("view_count + ?", 1),
			"utime":      time.Now().UnixMilli(),
		}),
	}).CreateInBatches(likes, len(ObjIDs))
	// todo 怎么更新缓存？
	return tx.Error
}

func (gad *GORMArticleLikeDao) InsertCollection(ctx context.Context, ObjType string, ObjID, uid int64) error {
	now := time.Now().UnixMilli()

	return gad.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.Clauses(clause.OnConflict{
			DoUpdates: clause.Assignments(map[string]any{ // 防止重复插入
				"status": 0,
				"utime":  time.Now().UnixMilli(),
			}),
		}).Create(&UserGiveCollect{
			ObjID:   ObjID,
			ObjType: ObjType,
			Uid:     uid,
			Ctime:   now,
			Utime:   now,
		}).Error
		if err != nil {
			return err
		}

		return tx.WithContext(ctx).Clauses(clause.OnConflict{
			DoUpdates: clause.Assignments(map[string]any{
				"collect_count": gorm.Expr("collect_count + ?", 1),
				"utime":         time.Now().UnixMilli(),
			}),
		}).Create(&Like{
			ObjID:        ObjID,
			ObjType:      ObjType,
			CollectCount: 1,
			Ctime:        now,
			Utime:        now,
		}).Error
	})
}

func (gad *GORMArticleLikeDao) DeleteLike(ctx context.Context, objType string, id int64, uid int64) error {
	now := time.Now().UnixMilli()
	return gad.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 此处where的顺序可以不用管，mysql会自动调整的
		err := tx.Model(&Like{}).Where("obj_id = ? and obj_type = ?", id, objType).
			Updates(map[string]any{
				"like_count": gorm.Expr("like_count - ?", 1),
				"utime":      now,
			}).Error
		if err != nil {
			return err
		}

		return tx.Model(&UserGiveLike{}).Where("obj_id = ? and obj_type = ? and uid = ?", id, objType, uid).
			Updates(map[string]any{
				"status": 1,
				"utime":  now,
			}).Error
	})
}

func (gad *GORMArticleLikeDao) InSertLike(ctx context.Context, objType string, id int64, uid int64) error {
	now := time.Now().UnixMilli()
	return gad.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 前端校验重复点赞，暂时此处不做校验了
		err := tx.Clauses(clause.OnConflict{
			DoUpdates: clause.Assignments(map[string]any{ // 防止重复插入
				"status": 0,
				"utime":  time.Now().UnixMilli(),
			}),
		}).Create(&UserGiveLike{ // 用户点赞表
			ObjID:   id,
			ObjType: objType,
			Uid:     uid,
			Ctime:   now,
			Utime:   now,
		}).Error
		if err != nil {
			return err
		}

		return tx.Clauses(clause.OnConflict{
			DoUpdates: clause.Assignments(map[string]any{
				"like_count": gorm.Expr("like_count + ?", 1),
				"utime":      time.Now().UnixMilli(),
			}),
		}).Create(&Like{
			ObjID:     id,
			ObjType:   objType,
			LikeCount: 1,
			Ctime:     now,
			Utime:     now,
		}).Error
	})
}

func (gad *GORMArticleLikeDao) IncrReadCount(ctx context.Context, ObjType string, ObjID int64) error {
	// 这两种分别是one和many，但不考虑insert的情况
	//return gad.db.WithContext(ctx).Where("obj_id = ? and obj_type = ?", ObjID, ObjType).
	//	Update("view_count", gorm.Expr("view_count + ?", 1)).Error

	//return gad.db.WithContext(ctx).Where("obj_id = ? and obj_type = ?", ObjID, ObjType).
	//	Updates(map[string]any{
	//		"view_count": gorm.Expr("view_count + ?", 1),
	//		"utime":      time.Now().UnixMilli(),
	//	}).Error

	// upsert
	now := time.Now().UnixMilli()
	return gad.db.WithContext(ctx).Clauses(clause.OnConflict{
		DoUpdates: clause.Assignments(map[string]any{
			"view_count": gorm.Expr("view_count + ?", 1),
			"utime":      time.Now().UnixMilli(),
		}),
		Columns: []clause.Column{{Name: "obj_id"}, {Name: "obj_type"}},
	}).Create(&Like{
		ObjID:     ObjID,
		ObjType:   ObjType,
		ViewCount: 1,
		Ctime:     now,
		Utime:     now,
	}).Error

	// 原生sql，upsert
	//sql := `INSERT INTO like (obj_id, obj_type, view_count) VALUES (?, ?, ?)
	//    ON DUPLICATE KEY UPDATE view_count = view_count + 1`
	//return gad.db.WithContext(ctx).Exec(sql, ObjID, ObjType, 1).Error
}

func (gad *GORMArticleLikeDao) GetPublishedByBatch(ctx context.Context, ObjType string, offset, limit int, now int64) ([]Like, error) {
	var res []Like
	return res, gad.db.WithContext(ctx).Where("obj_type = ? and ctime < ?", ObjType, now).
		Offset(offset).Limit(limit).Find(&res).Error
}

// TODO 找TOPn，超大规模数据情况下，归并

// Like 收集点赞数TOPn的数据
// 帖子的点赞，收藏，观看数
type Like struct {
	ID int64 `gorm:"primaryKey,autoIncrement"`

	// 联合索引， ObjID区分度更高，放左侧
	ObjID   int64  `gorm:"uniqueIndex:idx_objid_objtype"`
	ObjType string `gorm:"uniqueIndex:idx_objid_objtype;type:varchar(64)"`

	LikeCount    int64 `gorm:"column:like_count"`
	CollectCount int64 `gorm:"column:collect_count"`
	ViewCount    int64 `gorm:"column:view_count"`

	Ctime int64 `gorm:"index:idx_ctime"`
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

type UserGiveCollect struct {
	ID int64 `gorm:"primaryKey,autoIncrement"`

	// 此处的频繁查询应该是：查用户对帖子是否点赞，即where uid=? and objid=? and objtype=?
	// 考虑用户查看自己点赞的帖子，所以放左侧
	Uid     int64  `gorm:"uniqueIndex:idx_uid_objid_objtype"`
	ObjID   int64  `gorm:"uniqueIndex:idx_uid_objid_objtype"`
	ObjType string `gorm:"uniqueIndex:idx_uid_objid_objtype;type:varchar(64)"`

	Ctime int64
	Utime int64

	Status uint8 // 0：收藏，1：取消收藏
}

type UserGiveRead struct {
	ID int64 `gorm:"primaryKey,autoIncrement"`

	Uid     int64  `gorm:"uniqueIndex:idx_uid_objid_objtype"`
	ObjID   int64  `gorm:"uniqueIndex:idx_uid_objid_objtype"`
	ObjType string `gorm:"uniqueIndex:idx_uid_objid_objtype;type:varchar(64)"`

	Ctime int64
	Utime int64
}
