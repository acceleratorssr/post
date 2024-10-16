package gin_extra

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"strconv"
)

// todo 干掉包变量
var vector *prometheus.CounterVec

func InitCounter(opt prometheus.CounterOpts) {
	vector = prometheus.NewCounterVec(opt,
		[]string{"code"})
	prometheus.MustRegister(vector)
}

// WrapClaimsAndReq
// TODO 除此之外还可以考虑单独解析claims或者req，解决全部post
func WrapClaimsAndReq[Req any](fn func(*gin.Context, Req) (Response, error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req Req
		if err := ctx.Bind(&req); err != nil {
			err = fmt.Errorf("解析请求参数失败%w", err)
			FailWithMessage(ctx, Internal, err.Error())
			return
		}

		//claim, ok := ctx.Get("userClaims")
		//if !ok {
		//	ctx.AbortWithStatus(http.StatusUnauthorized)
		//	err := fmt.Errorf("无法获得 claims:%v", ctx.Request.URL.Path)
		//	FailWithMessage(ctx, Internal, err.Error())
		//	return
		//}
		//
		//claims, ok := claim.(*user.ClaimsUser)
		//if !ok {
		//	ctx.AbortWithStatus(http.StatusUnauthorized)
		//	err := fmt.Errorf("无法获得 claims:%v", ctx.Request.URL.Path)
		//	FailWithMessage(ctx, Internal, err.Error())
		//	return
		//}

		res, err := fn(ctx, req)

		if err != nil {
			err = fmt.Errorf("业务失败:%w", err)
			FailWithMessage(ctx, Internal, err.Error())
			return
		}

		OKWithDataAndMsg(ctx, res.Data, res.Msg)
	}
}

func WrapWithReq[Req any](fn func(*gin.Context, Req) (*Response, error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req Req
		if err := ctx.Bind(&req); err != nil {
			err = fmt.Errorf("解析请求参数失败%w", err)
			FailWithMessage(ctx, InvalidArgument, err.Error())
			return
		}

		res, err := fn(ctx, req)

		vector.WithLabelValues(strconv.Itoa(int(res.Code))).Inc()

		// _maxCode 为最大错误码
		if err != nil || (res.Code < _maxCode && res.Code > OK) {
			err = fmt.Errorf("业务失败:%w", err)
			FailWithMessage(ctx, res.Code, res.Msg+"-"+err.Error())
			return
		}

		OKWithDataAndMsg(ctx, res.Data, res.Msg)
	}
}

func WrapNilReq(fn func(*gin.Context) (*Response, error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		res, err := fn(ctx)

		vector.WithLabelValues(strconv.Itoa(int(res.Code))).Inc()

		if err != nil {
			err = fmt.Errorf("业务失败:%w", err)
			FailWithMessage(ctx, res.Code, err.Error())
			return
		}

		OKWithDataAndMsg(ctx, res.Data, res.Msg)
	}
}
