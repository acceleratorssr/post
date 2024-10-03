package domain

type User struct {
	UID        uint64
	Username   string
	Nickname   string
	Password   string
	TotpSecret string
	UserAgent  string
}
