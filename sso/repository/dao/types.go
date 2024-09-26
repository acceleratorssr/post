package dao

type User struct {
	ID         uint64 `gorm:"primaryKey,autoIncrement"`
	Username   string `gorm:"type:varchar(64);uniqueIndex"`
	Nickname   string `gorm:"type:varchar(64)"`
	Password   string `gorm:"type:varchar(64)"`
	TotpSecret string `gorm:"type:varchar(64)"`
	UserAgent  string `gorm:"type:varchar(255)"`
	QrcodeURL  string `gorm:"type:varchar(255)"`
}
