package domain

type User struct {
	Username   string
	Nickname   string
	Password   string
	TotpSecret string
	UserAgent  string
}
