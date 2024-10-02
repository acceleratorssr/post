package domain

type User struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Nickname string `json:"nickname"`
}

// UserInfo 可变用户信息
type UserInfo struct {
	Nickname string `json:"nickname"`
}
