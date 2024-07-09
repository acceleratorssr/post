package ioc

import (
	_ "embed"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"post/repository/dao"
)

//go:embed conf.yaml
var content string

func InitDB() *gorm.DB {
	type Config struct {
		DSN string `yaml:"dsn"`
	}
	c := Config{
		DSN: content,
	}

	db, err := gorm.Open(mysql.Open(c.DSN), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	err = dao.InitTables(db)
	if err != nil {
		panic(err)
	}
	return db
}
