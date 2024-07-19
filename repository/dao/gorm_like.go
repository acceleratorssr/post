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

func (gad *GORMArticleLikeDao) DeleteLike(ctx context.Context, objType string, id int64, uid int64) error {
	now := time.Now().UnixMilli()
	return gad.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 此处where的顺序可以不用管，mysql会自动调整的
		err := tx.Model(&Like{}).Where("obj_id = ? and obj_type = ? and uid = ?", id, objType, uid).
			Updates(map[string]any{
				"like_count": gorm.Expr("collect_count - ?", 1),
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
		//Columns: []clause.Column{{Name: "obj_id"}, {Name: "obj_type"}},
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
