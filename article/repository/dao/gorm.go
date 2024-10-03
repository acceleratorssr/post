package dao

import (
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

func (gad *GORMArticleDao) GetPublishedByID(ctx context.Context, id uint64) (*ArticleReader, error) {
	var art ArticleReader
	err := gad.db.WithContext(ctx).Where("id = ?", id).First(&art).Error
	return &art, err
}

func (gad *GORMArticleDao) GetByID(ctx context.Context, id uint64) (*ArticleAuthor, error) {
	var art ArticleAuthor
	err := gad.db.WithContext(ctx).Where("id = ?", id).First(&art).Error
	return &art, err
}

// GetListByAuthor 经典order加索引，此处authorid和utime可建立联合索引
// TODO sql性能优化
func (gad *GORMArticleDao) GetListByAuthor(ctx context.Context, uid uint64, limit int, offset int) ([]ArticleAuthor, error) {
	var arts []ArticleAuthor
	err := gad.db.WithContext(ctx).
		Where("authorid = ?", uid).
		Limit(limit).
		Offset(offset).
		//Order("utime DESC").
		Order(clause.OrderBy{Columns: []clause.OrderByColumn{
			{Column: clause.Column{Name: "utime"}, Desc: true},
		}}).
		Find(&arts).Error

	return arts, err
}

func (gad *GORMArticleDao) Insert(ctx context.Context, art *ArticleAuthor) (uint64, error) {
	art.Ctime = time.Now().UnixMilli()
	art.Utime = art.Ctime
	err := gad.db.WithContext(ctx).Create(art).Error
	return art.Id, err
}

func (gad *GORMArticleDao) InsertReader(ctx context.Context, art *ArticleReader) (uint64, error) {
	art.Ctime = time.Now().UnixMilli()
	art.Utime = art.Ctime
	err := gad.db.WithContext(ctx).Create(art).Error
	return art.Id, err
}

func (gad *GORMArticleDao) UpdateByID(ctx context.Context, art *ArticleAuthor) error {
	art.Utime = time.Now().UnixMilli()
	// 写简单，但是可读性不强，struct默认不会更新零值，map会
	res := gad.db.WithContext(ctx).Model(art).Where("authorid = ?", art.Authorid).Updates(*art)

	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("系统错误")
	}

	return res.Error
}

// SyncStatus
// todo 此处没有对like进行更新
func (gad *GORMArticleDao) SyncStatus(ctx context.Context, art *ArticleAuthor) error {
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
