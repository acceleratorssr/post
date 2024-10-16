package connpool

import (
	"github.com/stretchr/testify/require"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"sync/atomic"
	"testing"
)

type Like struct {
	ID int64 `gorm-extra:"primaryKey,autoIncrement"`

	// 联合索引， ObjID区分度更高，放左侧
	ObjID   int64  `gorm-extra:"uniqueIndex:idx_objid_objtype"`
	ObjType string `gorm-extra:"uniqueIndex:idx_objid_objtype;type:varchar(64)"`

	LikeCount    int64 `gorm-extra:"column:like_count"`
	CollectCount int64 `gorm-extra:"column:collect_count"`
	ViewCount    int64 `gorm-extra:"column:view_count"`

	Ctime int64 `gorm-extra:"index:idx_ctime"`
	Utime int64
}

func TestConnPool(t *testing.T) {
	src, _ := gorm.Open(
		mysql.Open("root:20031214pzw!@tcp(127.0.0.1:3306)/like_demo"))

	dst, _ := gorm.Open(
		mysql.Open("root:20031214pzw!@tcp(127.0.0.1:3306)/test"))

	var v atomic.Value
	v.Store(PattenDependBase)
	db, err := gorm.Open(mysql.New(mysql.Config{
		Conn: &DoubleWritePool{
			base:   src.ConnPool,
			target: dst.ConnPool,
			patten: v,
		},
	}))
	require.NoError(t, err)
	t.Log(db)
	db.Create(&Like{
		ObjID:        100,
		ObjType:      "article",
		LikeCount:    1,
		CollectCount: 1,
		ViewCount:    1,
	})

	db.Transaction(func(tx *gorm.DB) error {
		return tx.Create(&Like{
			ObjID:        101,
			ObjType:      "article",
			LikeCount:    1,
			CollectCount: 1,
			ViewCount:    1,
		}).Error
	})
}
