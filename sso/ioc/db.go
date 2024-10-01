package ioc

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"post/sso/config"
	"post/sso/repository/dao"
)

func InitDB(info *config.Info) *gorm.DB {
	DSN := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		info.Config.Mysql.User,
		info.Config.Mysql.Password,
		info.Config.Mysql.Host,
		info.Config.Mysql.Port,
		info.Config.Mysql.DB,
	)
	db, err := gorm.Open(mysql.Open(DSN), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	err = dao.InitTables(db)
	if err != nil {
		panic(err)
	}

	//sqlDB, err := db.DB()
	//
	//sqlDB.SetMaxIdleConns(10)
	//
	//sqlDB.SetMaxOpenConns(100)
	//
	//sqlDB.SetConnMaxLifetime(time.Hour)

	return db
}
