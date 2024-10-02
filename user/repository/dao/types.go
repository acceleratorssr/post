package dao

type User struct {
	ID       int64  `gorm:"primaryKey,autoIncrement"`
	Username string `gorm:"type:varchar(64);uniqueIndex"`
	Nickname string `gorm:"type:varchar(64)"`

	Utime int64
	Ctime int64
}

// UserInfo 可变用户信息
type UserInfo struct {
	Nickname string
}
