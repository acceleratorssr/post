package domain

import (
	"context"
	"github.com/golang-jwt/jwt/v5"
)

type JwtPayload struct {
	UID      uint64
	Username string
	NickName string
}

type Claims struct {
	*JwtPayload
	jwt.RegisteredClaims
}

type AuthService interface {
	GenerateRefreshToken(ctx context.Context, user *JwtPayload) (string, error)
	GenerateAccessToken(ctx context.Context, user *JwtPayload) (string, error)
	ValidateToken(ctx context.Context, token string) (*Claims, error)
	GetPublicKey(ctx context.Context) string
}
