package ioc

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"post/user/config"
	"post/user/repository/dao"
)

func InitDB() *gorm.DB {
	info := config.InitConfig()

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
	return db
}
