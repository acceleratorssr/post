package gin_ex

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
	"post/internal/domain"
	"post/internal/user"
	"post/internal/utils"
	"strconv"
)

var vector *prometheus.CounterVec

func InitCounter(opt prometheus.CounterOpts) {
	vector = prometheus.NewCounterVec(opt,
		[]string{"code"})
	prometheus.MustRegister(vector)
}

// WrapClaimsAndReq
// TODO 除此之外还可以考虑单独解析claims或者req，解决全部post
func WrapClaimsAndReq[Req any](fn func(context.Context, Req, user.ClaimsUser) (utils.Response, error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req Req
		if err := ctx.Bind(&req); err != nil {
			err = fmt.Errorf("解析请求失败%w", err)
			utils.FailWithMessage(domain.ErrSystem, err.Error(), ctx)
			return
		}

		claim, ok := ctx.Get("userClaims")
		if !ok {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			err := fmt.Errorf("无法获得 claims:%v", ctx.Request.URL.Path)
			utils.FailWithMessage(domain.ErrSystem, err.Error(), ctx)
			return
		}

		claims, ok := claim.(user.ClaimsUser)
		if !ok {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			err := fmt.Errorf("无法获得 claims:%v", ctx.Request.URL.Path)
			utils.FailWithMessage(domain.ErrSystem, err.Error(), ctx)
			return
		}

		res, err := fn(ctx.Request.Context(), req, claims)

		if err != nil {
			err = fmt.Errorf("业务失败:%w", err)
			utils.FailWithMessage(domain.ErrSystem, err.Error(), ctx)
		}

		utils.OK(res.Data, res.Msg, ctx)
	}
}

func WrapNilReq(fn func(context.Context) (utils.Response, error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		res, err := fn(ctx.Request.Context())

		vector.WithLabelValues(strconv.Itoa(res.Code)).Inc()

		// utils.UserInvalidInput 为最小错误码
		if err != nil || res.Code >= utils.UserInvalidInput {
			err = fmt.Errorf("业务失败:%w", err)
			utils.FailWithMessage(domain.StatusType(res.Code), err.Error(), ctx)
		}

		utils.OK(res.Data, res.Msg, ctx)
	}
}
