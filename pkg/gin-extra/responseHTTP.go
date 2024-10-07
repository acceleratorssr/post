package gin_ex

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
	"reflect"
)

type Response struct {
	Code Code   `json:"code"`
	Data any    `json:"data"`
	Msg  string `json:"msg"`
}

type ListResponse[T any] struct {
	TotalCount int `json:"count"`
	List       T   `json:"list"`
}

func Result(ctx *gin.Context, code Code, data any, msg string) {
	ctx.JSON(http.StatusOK, Response{
		Code: code,
		Data: data,
		Msg:  msg,
	})
}

func OKWithDataAndMsg(ctx *gin.Context, data any, msg string) {
	Result(ctx, http.StatusOK, data, msg)
}

func OKWithData(ctx *gin.Context, data any) {
	Result(ctx, http.StatusOK, data, "")
}

func OKWithMessage(ctx *gin.Context, msg string) {
	Result(ctx, http.StatusOK, nil, msg)
}

func OKWithList[T any](ctx *gin.Context, list T, totalCount int) {
	OKWithData(ctx, ListResponse[T]{
		List:       list,
		TotalCount: totalCount,
	})
}

func Fail(ctx *gin.Context, code Code, data any, msg string) {
	Result(ctx, code, data, msg)
}

func FailWithCode(ctx *gin.Context, code Code) {
	msg := canonicalString(code)
	Result(ctx, code, nil, msg)
}

func FailWithMessage(ctx *gin.Context, code Code, msg string) {
	Result(ctx, code, nil, msg)
}

// FailWithError 从验证错误中提取字段的自定义错误消息，obj为对应结构体
func FailWithError(ctx *gin.Context, err error, obj any) {
	msg := GetValidMsg(err, obj)
	FailWithMessage(ctx, Unknown, msg)
}

func GetValidMsg(err error, obj any) string {
	getObj := reflect.TypeOf(obj)
	var errs validator.ValidationErrors

	if errors.As(err, &errs) {
		for _, e := range errs {
			if f, exits := getObj.FieldByName(e.Field()); exits {
				//如果字段存在，通过 Tag.GetListInfo("msg") 获取字段的自定义错误消息，并将其作为函数的返回值。
				return f.Tag.Get("msg")
			}
		}
	}

	return err.Error()
}
