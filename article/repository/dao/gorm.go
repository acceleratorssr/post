package dao

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

// GORMArticleDao todo 考虑分库
type GORMArticleDao struct {
	db *gorm.DB
}

func (gad *GORMArticleDao) ListByID(ctx context.Context, uid uint64, list *List) ([]ArticleReader, error) {
	var arts []ArticleReader
	err := gad.db.WithContext(ctx).Where(clause.Expr{SQL: fmt.Sprintf("%s < ?", list.OrderBy), Vars: []interface{}{list.LastValue}}).
		Where("authorid = ?", uid).Omit("content", "id").Limit(list.Limit).
		Order(clause.OrderByColumn{Column: clause.Column{Name: list.OrderBy}, Desc: list.Desc}).Find(&arts).Error

	return arts, err
}

func (gad *GORMArticleDao) ListPublished(ctx context.Context, list *List) ([]ArticleReader, error) {
	var arts []ArticleReader
	// sql 优化
	err := gad.db.WithContext(ctx).Where(clause.Expr{SQL: fmt.Sprintf("%s < ?", list.OrderBy), Vars: []interface{}{list.LastValue}}).
		Limit(list.Limit).Omit("id").
		Order(clause.OrderByColumn{Column: clause.Column{Name: list.OrderBy}, Desc: list.Desc}).Find(&arts).Error // todo 到源码里看看

	// 错误生成sql：
	//err := gad.db.WithContext(ctx).Omit("content", "id").
	//	Offset(int(list.LastValue)).Limit(list.Limit).
	//	Order(clause.OrderBy{Columns: []clause.OrderByColumn{
	//		{Column: clause.Column{Name: list.OrderBy}, Desc: list.Desc},
	//	}}).Find(&arts).Error

	// 未优化sql：
	//err := gad.db.WithContext(ctx).Omit("content", "id").
	//	Offset(int(list.LastValue)).Limit(list.Limit).
	//	Order("ctime desc").Find(&arts).Error

	return arts, err
}

func (gad *GORMArticleDao) GetPublishedByIDs(ctx context.Context, aids []uint64) ([]ArticleReader, error) {
	var articles []ArticleReader
	err := gad.db.WithContext(ctx).Where("snow_id IN ?", aids).Find(&articles).Error
	if err != nil {
		return nil, err
	}
	return articles, nil
}

func (gad *GORMArticleDao) GetAuthorByID(ctx context.Context, aid, uid uint64) (*ArticleAuthor, error) {
	var art ArticleAuthor
	err := gad.db.WithContext(ctx).Where("snow_id = ?", aid).First(&art).Error
	if err == nil && art.Authorid != uid {
		// 监控，uid用户非法查询他人文章
		return nil, gorm.ErrRecordNotFound
	}
	return &art, err
}

func (gad *GORMArticleDao) GetPublishedByID(ctx context.Context, aid uint64) (*ArticleReader, error) {
	var art ArticleReader
	err := gad.db.WithContext(ctx).Where("snow_id = ?", aid).First(&art).Error
	return &art, err
}

// GetListByAuthor 经典order加索引，此处authorid和utime可建立联合索引
// TODO sql性能优化
func (gad *GORMArticleDao) GetListByAuthor(ctx context.Context, uid uint64, list *List) ([]ArticleAuthor, error) {
	var arts []ArticleAuthor
	err := gad.db.WithContext(ctx).Where(clause.Expr{SQL: fmt.Sprintf("%s < ?", list.OrderBy), Vars: []interface{}{list.LastValue}}).
		Where("authorid = ?", uid).Omit("content", "id").Limit(list.Limit).
		Order(clause.OrderByColumn{Column: clause.Column{Name: list.OrderBy}, Desc: list.Desc}).Find(&arts).Error

	return arts, err
}

func (gad *GORMArticleDao) Insert(ctx context.Context, art *ArticleAuthor) error {
	art.Ctime = time.Now().UnixMilli()
	art.Utime = art.Ctime
	return gad.db.WithContext(ctx).Create(art).Error
}

func (gad *GORMArticleDao) UpsertReader(ctx context.Context, art *ArticleReader) error {
	art.Ctime = time.Now().UnixMilli()
	art.Utime = art.Ctime

	return gad.db.WithContext(ctx).Clauses(
		clause.OnConflict{
			Columns:   []clause.Column{{Name: "snow_id"}},
			DoUpdates: clause.AssignmentColumns([]string{"title", "content", "utime"}),
		},
	).Create(art).Error
}

func (gad *GORMArticleDao) UpdateByID(ctx context.Context, art *ArticleAuthor) error {
	art.Utime = time.Now().UnixMilli()
	// 写简单，但是可读性不强，struct默认不会更新零值，map会
	res := gad.db.WithContext(ctx).Model(art).Where("snow_id = ?", art.SnowID).Updates(*art)

	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (gad *GORMArticleDao) DeleteReader(ctx context.Context, aid uint64, uid uint64) error {
	res := gad.db.WithContext(ctx).Where("snow_id = ? AND author_id = ?", aid, uid).Delete(&ArticleReader{})
	if res.Error != nil {
		// 数据库有问题
		return res.Error
	}
	if res.RowsAffected == 0 {
		// 监控，uid用户非法删除他人文章
		return gorm.ErrRecordNotFound
	}
	return nil
}

func NewGORMArticleDao(db *gorm.DB) ArticleDao {
	return &GORMArticleDao{
		db: db,
	}
}
