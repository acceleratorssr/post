package domain

import "github.com/golang-jwt/jwt/v5"

type JwtPayload struct {
	Username    string
	UserID      uint64
	NickName    string
	Permissions int
}

type Claims struct {
	*JwtPayload
	jwt.RegisteredClaims
}

type User struct {
	ID          uint64
	Username    string
	Nickname    string
	Password    string
	TotpSecret  string
	UserAgent   string
	QrcodeURL   string
	Permissions int
}
