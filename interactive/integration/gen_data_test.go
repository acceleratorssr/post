package integration

import (
	_ "embed"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
	"math/rand"
	"os"
	"post/interactive/integration/startup"
	"post/interactive/repository/dao"
	"strconv"
	"testing"
	"time"
)

//go:embed init.sql
var initSQL string

// TestGenSQL 这个测试就是用来生成 SQL
func TestGenSQL(t *testing.T) {
	file, err := os.OpenFile("data.sql",
		os.O_RDWR|os.O_APPEND|os.O_CREATE|os.O_TRUNC, 0666)
	require.NoError(t, err)
	defer file.Close()

	// 创建数据库和数据表的语句，防止还没初始化
	_, err = file.WriteString(initSQL)
	require.NoError(t, err)

	const prefix = "INSERT INTO `likes`(`obj_id`, `obj_type`, `view_count`, `collect_count`, `like_count`, `ctime`, `utime`)\nVALUES"
	const rowNum = 10

	now := time.Now().UnixMilli()
	_, err = file.WriteString(prefix)

	for i := 0; i < rowNum; i++ {
		if i > 0 {
			file.Write([]byte{',', '\n'})
		}
		file.Write([]byte{'('})
		// biz_id
		file.WriteString(strconv.Itoa(i + 1))
		// biz
		file.WriteString(`,"test",`)
		// read_cnt
		file.WriteString(strconv.Itoa(int(rand.Int31n(10000))))
		file.Write([]byte{','})

		// collect_cnt
		file.WriteString(strconv.Itoa(int(rand.Int31n(10000))))
		file.Write([]byte{','})
		// like_cnt
		file.WriteString(strconv.Itoa(int(rand.Int31n(10000))))
		file.Write([]byte{','})

		// ctime
		file.WriteString(strconv.FormatInt(now, 10))
		file.Write([]byte{','})

		// utime
		file.WriteString(strconv.FormatInt(now, 10))

		file.Write([]byte{')'})
	}
}

func TestGenData(t *testing.T) {
	// 批量插入，数据量不是特别大的时候
	// GenData 要比 GenSQL 慢
	// 根据需要调整批次，和每个批次大小
	db := startup.InitDB()

	// 模拟执行 SQL 查询的选项
	// 为 true 时，GORM 不会实际执行查询
	// 而是生成并返回将要执行的 SQL 语句
	db.DryRun = false

	// 1000 批
	for i := 0; i < 10; i++ {
		// 每次 100 条
		// 可以考虑直接用 CreateInBatches，GORM 帮你分批次
		const batchSize = 100
		data := make([]dao.Like, 0, batchSize)
		now := time.Now().UnixMilli()
		for j := 0; j < batchSize; j++ {
			data = append(data, dao.Like{
				ObjType:      "test",
				ObjID:        int64(i*batchSize + j + 1),
				ViewCount:    rand.Int63(),
				LikeCount:    rand.Int63(),
				CollectCount: rand.Int63(),
				Utime:        now,
				Ctime:        now,
			})
		}

		err := db.Transaction(func(tx *gorm.DB) error {
			return tx.Create(data).Error
		})
		require.NoError(t, err)
	}
}
