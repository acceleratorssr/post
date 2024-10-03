package dao

type User struct {
	ID       uint64 `gorm:"primaryKey"`
	Username string `gorm:"type:varchar(64);uniqueIndex"`
	Nickname string `gorm:"type:varchar(64)"`

	Utime int64
	Ctime int64
}

// UserInfo 可变用户信息
type UserInfo struct {
	Nickname string
}
