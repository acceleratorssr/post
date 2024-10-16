package middleware

import (
	"context"
	"crypto/ecdsa"
	_ "embed"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"post/pkg/gin-extra"
	"post/sso/domain"
	"strings"
	"time"
)

// Jwt 如果需要保证立刻登出的效果，可结合session
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

	if claims, ok := token.Claims.(*domain.Claims); ok && token.Valid && !claims.ExpiresAt.Time.Before(time.Now()) {
		if issuer, _ := token.Claims.GetIssuer(); issuer == "SSO" {
			return claims, nil
		}
		return nil, jwt.ErrTokenInvalidIssuer
	}
	return nil, jwt.ErrTokenUnverifiable
}

func (j *Jwt) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// X-Refresh-Token: <refresh_token>
		xToken := ctx.GetHeader("X-Refresh-Token")
		// Authorization: Bearer <token>
		token := ctx.GetHeader("Authorization")
		tokens := strings.Split(token, " ")

		if xToken != "" { // 仅用于logout
			xClaims, err := j.validateToken(ctx, xToken)
			if err != nil {
				// 部署到k8s内时，直接json格式日志到stdout即可；
				gin_extra.FailWithMessage(ctx, gin_extra.Unauthenticated, "请重新登录")
				ctx.AbortWithStatus(http.StatusUnauthorized)
				return
			}
			ctx.Set("x_exp", xClaims.ExpiresAt.Unix())
			ctx.Set("x_token", xToken)
		}

		if len(tokens) != 2 {
			if xToken == "" {
				gin_extra.FailWithMessage(ctx, gin_extra.InvalidArgument, "token错误")
				ctx.AbortWithStatus(http.StatusUnauthorized)
				return
			}
		} else { // 放行单refresh token
			claims, err := j.validateToken(ctx, tokens[1])
			if err != nil {
				// 重定向到refresh，再回来
				ctx.AbortWithStatus(http.StatusUnauthorized)
				return
			}

			// userClaims 内的数据 ---
			ctx.Set("uid", claims.UID)
			ctx.Set("username", claims.Username)
			ctx.Set("nickname", claims.NickName)
			ctx.Set("exp", claims.ExpiresAt.Unix())
			// ---

			ctx.Set("token", tokens[1]) // 懒得二次处理，所以此处暂时放到ctx中
		}

		ctx.Next()
	}
}
