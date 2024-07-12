package dao

import (
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type GORMArticleDao struct {
	db *gorm.DB
}

func NewGORMArticleDao(db *gorm.DB) ArticleDao {
	return &GORMArticleDao{
		db: db,
	}
}

func (gad *GORMArticleDao) Insert(ctx context.Context, art ArticleAuthor) (int64, error) {
	art.Ctime = time.Now().UnixMilli()
	art.Utime = art.Ctime
	err := gad.db.WithContext(ctx).Create(&art).Error
	return art.Id, err
}

func (gad *GORMArticleDao) UpdateByID(ctx context.Context, art ArticleAuthor) error {
	art.Utime = time.Now().UnixMilli()
	// 写简单，但是可读性不强，从代码看不出哪些字段不是零值，会被更新
	res := gad.db.WithContext(ctx).Model(&art).Where("authorid = ?", art.Authorid).Updates(art)

	// 可读性强，但是如果后面字段要更新，这里需要修改
	//err := gad.db.WithContext(ctx).Model(&art).Where("id=?", art.Id).
	//	Updates(map[string]any{
	//		"title":   art.Title,
	//		"content": art.Content,
	//		"utime":   art.Utime,
	//	}).Error

	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("系统错误")
	}

	return res.Error
}

func (gad *GORMArticleDao) SyncStatus(ctx context.Context, art ArticleAuthor) error {
	now := time.Now().UnixMilli()
	return gad.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		res := tx.Model(&ArticleReader{}).Where("id=?", art.Id).Delete(&ArticleReader{})
		if res.Error != nil {
			// 数据库有问题
			return res.Error
		}
		if res.RowsAffected == 0 {
			// 找不到记录，id错误
			return errors.New("id不匹配")
		}

		return tx.Model(&ArticleAuthor{}).Where("id=?", art.Id).Updates(map[string]any{
			"status": art.Status,
			"utime":  now,
		}).Error
	})
}
