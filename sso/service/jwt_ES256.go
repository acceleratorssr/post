package service

import (
	"context"
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"github.com/golang-jwt/jwt/v5"
	"os"
	"post/sso/config"
	"post/sso/domain"
	"time"
)

type authService struct {
	info *config.Info
}

func (a *authService) GetPublicKey(ctx context.Context) string {
	ecdsaPubKey := a.loadPublicKey()
	pubKeyBytes, err := x509.MarshalPKIXPublicKey(ecdsaPubKey)
	if err != nil {
		return ""
	}
	return base64.StdEncoding.EncodeToString(pubKeyBytes)
}

func (a *authService) GenerateRefreshToken(ctx context.Context, user *domain.JwtPayload) (string, error) {
	claims := &domain.Claims{
		JwtPayload: user,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(a.info.Config.Jwt.LongExpires) * time.Hour)),
			Issuer:    a.info.Config.Jwt.Issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	return token.SignedString(a.loadPrivateKey())
}

func (a *authService) GenerateAccessToken(ctx context.Context, user *domain.JwtPayload) (string, error) {
	claims := &domain.Claims{
		JwtPayload: user,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(a.info.Config.Jwt.LongExpires) * time.Hour)),
			Issuer:    a.info.Config.Jwt.Issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	return token.SignedString(a.loadPrivateKey())
}

func (a *authService) ValidateToken(ctx context.Context, tokenStr string) (*domain.Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &domain.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return a.loadPublicKey(), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*domain.Claims); ok && token.Valid {
		if issuer, _ := token.Claims.GetIssuer(); issuer == "SSO" {
			return claims, nil
		}
		return nil, jwt.ErrTokenInvalidIssuer
	}
	return nil, jwt.ErrTokenUnverifiable
}

func (a *authService) loadPrivateKey() *ecdsa.PrivateKey {
	//privPEM, err := os.ReadFile("../config/private_key.pem")
	privPEM, err := os.ReadFile("./sso/config/private_key.pem")
	if err != nil {
		return nil
	}

	block, _ := pem.Decode(privPEM)
	if block == nil || block.Type != "EC PRIVATE KEY" {
		return nil
	}

	privKey, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		return nil
	}

	return privKey
}

func (a *authService) loadPublicKey() *ecdsa.PublicKey {
	//pubPEM, err := os.ReadFile("../config/public_key.pem")
	pubPEM, err := os.ReadFile("./sso/config/public_key.pem")
	if err != nil {
		return nil
	}

	block, _ := pem.Decode(pubPEM)
	if block == nil || block.Type != "PUBLIC KEY" {
		return nil
	}

	pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil
	}

	ecdsaPubKey, ok := pubKey.(*ecdsa.PublicKey)
	if !ok {
		return nil
	}

	return ecdsaPubKey
}

func NewJWTService(info *config.Info) domain.AuthService {
	return &authService{
		info: info,
	}
}
