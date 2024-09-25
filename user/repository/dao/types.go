package dao

type User struct {
	ID          int64  `gorm:"primaryKey,autoIncrement"`
	Username    string `gorm:"type:varchar(64);uniqueIndex"`
	Password    string `gorm:"type:varchar(64)"`
	Nickname    string `gorm:"type:varchar(64)"`
	Permissions int
}
