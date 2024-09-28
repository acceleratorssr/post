package domain

import "github.com/golang-jwt/jwt/v5"

type JwtPayload struct {
	Username string
	NickName string
}

type Claims struct {
	*JwtPayload
	jwt.RegisteredClaims
}

type User struct {
	Username   string
	Nickname   string
	Password   string
	TotpSecret string
	UserAgent  string
}
