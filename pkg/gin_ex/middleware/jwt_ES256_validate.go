package middleware

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"post/sso/domain"
	"strings"
)

type Jwt struct {
	publicKey *ecdsa.PublicKey
}

func NewJwt(key *ecdsa.PublicKey) *Jwt {
	return &Jwt{
		publicKey: key,
	}
}

func (j *Jwt) validateToken(ctx context.Context, tokenStr string) (*domain.Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &domain.Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
			panic(fmt.Errorf("unexpected signing method: %v", token.Header["alg"]))
		}
		return j.publicKey, nil
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

func (j *Jwt) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Authorization: Bearer <token>
		token := ctx.GetHeader("Authorization")
		tokens := strings.Split(token, " ")
		if len(tokens) != 2 {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		claims, err := j.validateToken(ctx, tokens[1])
		if err != nil {
			// 部署到k8s内时，直接json格式日志到stdout即可；
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		ctx.Set("username", claims.Username)
		ctx.Set("nickname", claims.NickName)
		ctx.Set("token", tokens[1]) // 懒得二次处理，所以此处暂时放到ctx中
		//ctx.Set("userAgent", ctx.Request.UserAgent())

		ctx.Next()
	}
}
