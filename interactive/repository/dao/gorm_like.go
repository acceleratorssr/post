package dao

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strings"
	"time"
)

type GORMArticleLikeDao struct {
	db *gorm.DB
}

// UpdateReadCountMany todo 优化索引
func (gad *GORMArticleLikeDao) UpdateReadCountMany(ctx context.Context, objType string, hmap map[uint64]int64) error {
	var ids []uint64
	caseStatements := make([]string, 0, len(hmap))
	insertParams := make([]interface{}, 0, len(hmap)*3) // 每条记录三个值：obj_id, obj_type, view_count
	updateParams := make([]interface{}, 0, len(hmap)*2) // 更新语句的参数：obj_id, view_count

	for id, count := range hmap {
		ids = append(ids, id)
		caseStatements = append(caseStatements, "WHEN obj_id = ? THEN view_count + ?")
		updateParams = append(updateParams, id, count)
		insertParams = append(insertParams, id, objType, count) // 插入的参数
	}

	// 生成 CASE 语句
	caseSQL := fmt.Sprintf("CASE %s END", strings.Join(caseStatements, " "))

	// 生成批量插入的占位符
	valuesPlaceholder := generateValuesPlaceholder(len(hmap))

	// 拼接最终的 SQL 语句
	upsertSQL := fmt.Sprintf(`
    INSERT INTO likes (obj_id, obj_type, view_count)
    VALUES %s
    ON DUPLICATE KEY UPDATE view_count = %s`,
		valuesPlaceholder,
		caseSQL,
	)

	// 将更新参数与插入参数合并
	finalParams := append(insertParams, updateParams...)

	return gad.db.WithContext(ctx).
		Exec(upsertSQL, finalParams...).Error
}

func (gad *GORMArticleLikeDao) IncrReadCountMany(ctx context.Context, objType string, objIDs []uint64) error {
	now := time.Now().UnixMilli()

	likeMap := make(map[uint64]*Like)
	for _, objID := range objIDs {
		if like, exists := likeMap[objID]; exists {
			like.ViewCount++
			like.Utime = now
		} else {
			likeMap[objID] = &Like{
				ObjID:     objID,
				ObjType:   objType,
				ViewCount: 1,
				Ctime:     now,
				Utime:     now,
			}
		}
	}

	likes := make([]Like, 0, len(likeMap))
	for _, like := range likeMap {
		likes = append(likes, *like)
	}

	tx := gad.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "obj_id"}, {Name: "obj_type"}}, // 指定冲突列
		DoUpdates: clause.Assignments(map[string]any{
			"view_count": gorm.Expr("view_count + ?", 1),
			"utime":      time.Now().UnixMilli(),
		}),
	}).CreateInBatches(likes, len(objIDs))
	// todo 怎么更新缓存？
	return tx.Error
}

func (gad *GORMArticleLikeDao) InsertCollection(ctx context.Context, objType string, objID, uid uint64) error {
	now := time.Now().UnixMilli()

	return gad.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.Clauses(clause.OnConflict{
			DoUpdates: clause.Assignments(map[string]any{ // 防止重复插入
				"status": 0,
				"utime":  time.Now().UnixMilli(),
			}),
		}).Create(&UserGiveCollect{
			ObjID:   objID,
			ObjType: objType,
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
			ObjID:        objID,
			ObjType:      objType,
			CollectCount: 1,
			Ctime:        now,
			Utime:        now,
		}).Error
	})
}

func (gad *GORMArticleLikeDao) DeleteLike(ctx context.Context, objType string, objID, uid uint64) error {
	now := time.Now().UnixMilli()
	return gad.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 此处where的顺序可以不用管，mysql会自动调整的
		err := tx.Model(&UserGiveLike{}).Where("obj_id = ? and obj_type = ? and uid = ?", objID, objType, uid).
			Updates(map[string]any{
				"status": 1,
				"utime":  now,
			}).Error
		if err != nil {
			return err
		}

		// 防止重复 取消 点赞，如果上述用户行为发生冲突且没有更新（unlike -> like）
		// 则视为重复 取消 点赞，直接终止事务
		rowsAffected := tx.RowsAffected
		if rowsAffected == 0 {
			return fmt.Errorf("重复点赞，操作终止")
		}

		return tx.Model(&Like{}).Where("obj_id = ? and obj_type = ?", objID, objType).
			Updates(map[string]any{
				"like_count": gorm.Expr("like_count - ?", 1),
				"utime":      now,
			}).Error
	})
}

func (gad *GORMArticleLikeDao) InSertLike(ctx context.Context, objType string, objID, uid uint64) error {
	now := time.Now().UnixMilli()
	return gad.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.Clauses(clause.OnConflict{
			DoUpdates: clause.Assignments(map[string]any{ // 防止重复插入
				"status": 0,
				"utime":  time.Now().UnixMilli(),
			}),
		}).Create(&UserGiveLike{ // 用户点赞表
			ObjID:   objID,
			ObjType: objType,
			Uid:     uid,
			Ctime:   now,
			Utime:   now,
		}).Error
		if err != nil {
			return err
		}

		// 防止重复点赞，如果上述用户行为发生冲突且没有更新（unlike -> like）
		// 则视为重复点赞，直接终止事务
		rowsAffected := tx.RowsAffected
		if rowsAffected == 0 {
			return fmt.Errorf("重复点赞，操作终止")
		}

		return tx.Clauses(clause.OnConflict{
			DoUpdates: clause.Assignments(map[string]any{
				"like_count": gorm.Expr("like_count + ?", 1),
				"utime":      time.Now().UnixMilli(),
			}),
		}).Create(&Like{
			ObjID:     objID,
			ObjType:   objType,
			LikeCount: 1,
			Ctime:     now,
			Utime:     now,
		}).Error
	})
}

func (gad *GORMArticleLikeDao) IncrReadCount(ctx context.Context, objType string, objID uint64) error {
	// 这两种分别是one和many，但不考虑insert的情况
	//return gad.db.WithContext(ctx).Where("obj_id = ? and obj_type = ?", ObjID, ObjType).
	//	Update("view_count", gorm-extra.Expr("view_count + ?", 1)).Error

	//return gad.db.WithContext(ctx).Where("obj_id = ? and obj_type = ?", ObjID, ObjType).
	//	Updates(map[string]any{
	//		"view_count": gorm-extra.Expr("view_count + ?", 1),
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
		ObjID:     objID,
		ObjType:   objType,
		ViewCount: 1,
		Ctime:     now,
		Utime:     now,
	}).Error

	// 原生sql，upsert
	//sql := `INSERT INTO like (obj_id, obj_type, view_count) VALUES (?, ?, ?)
	//    ON DUPLICATE KEY UPDATE view_count = view_count + 1`
	//return gad.db.WithContext(ctx).Exec(sql, ObjID, ObjType, 1).Error
}

func (gad *GORMArticleLikeDao) GetLikeByBatch(ctx context.Context, objType string, limit int, lastValue int64, orderBy string, desc bool) ([]Like, error) {
	var res []Like
	err := gad.db.WithContext(ctx).Where(clause.Expr{SQL: fmt.Sprintf("%s < ?", orderBy), Vars: []interface{}{lastValue}}).
		Where("obj_type = ?", objType).Limit(limit).
		Order(clause.OrderByColumn{Column: clause.Column{Name: orderBy}, Desc: desc}).Find(&res).Error

	return res, err
}

func generateValuesPlaceholder(n int) string {
	placeholders := make([]string, n)
	for i := range placeholders {
		placeholders[i] = "(?, ?, ?)"
	}
	return strings.Join(placeholders, ", ")
}

func NewGORMArticleLikeDao(db *gorm.DB) ArticleLikeDao {
	return &GORMArticleLikeDao{
		db: db,
	}
}
