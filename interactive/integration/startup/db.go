package startup

import (
	"context"
	"database/sql"
	_ "embed"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"post/interactive/repository/dao"
	"time"
)

var db *gorm.DB

//go:embed mysql.yaml
var mysqlDSN string

func InitDB() *gorm.DB {
	if db == nil {
		sqlDB, err := sql.Open("mysql", mysqlDSN)
		if err != nil {
			panic(err)
		}
		for {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			err = sqlDB.PingContext(ctx)
			cancel()
			if err == nil {
				break
			}
		}
		db, err = gorm.Open(mysql.Open(mysqlDSN))
		if err != nil {
			panic(err)
		}
		err = dao.InitTables(db)
		if err != nil {
			panic(err)
		}
		//db = db.Debug()
	}
	return db
}
