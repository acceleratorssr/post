package dao

type User struct {
	ID          int64  `gorm:"primaryKey,autoIncrement"`
	Username    string `gorm:"type:varchar(64);uniqueIndex"`
	Nickname    string `gorm:"type:varchar(64)"`
	Permissions int
}
