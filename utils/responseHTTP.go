package utils

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
	"post/domain"
	"reflect"
)

type Response struct {
	Code int    `json:"code"`
	Data any    `json:"data"`
	Msg  string `json:"msg"`
}

type ListResponse[T any] struct {
	TotalCount int `json:"count"`
	List       T   `json:"list"`
}

func Result(code int, data any, msg string, c *gin.Context) {
	c.JSON(http.StatusOK, Response{
		Code: code,
		Data: data,
		Msg:  msg,
	})
}

func OK(data any, msg string, c *gin.Context) {
	Result(http.StatusOK, data, msg, c)
}

func OKWithData(data any, c *gin.Context) {
	Result(http.StatusOK, data, "", c)
}

func OKWithMessage(msg string, c *gin.Context) {
	Result(http.StatusOK, nil, msg, c)
}

func OKWithList[T any](list T, totalCount int, c *gin.Context) {
	OKWithData(ListResponse[T]{
		List:       list,
		TotalCount: totalCount,
	}, c)
}

func Fail(code domain.StatusType, data any, msg string, c *gin.Context) {
	Result(code.ToInt(), data, msg, c)
}

func FailWithCode(code domain.StatusType, c *gin.Context) {
	msg, ok := domain.ErrorMap[code]
	if !ok {
		msg = "Unknown error"
	}
	Result(code.ToInt(), nil, msg, c)
}

func FailWithMessage(code domain.StatusType, msg string, c *gin.Context) {
	Result(code.ToInt(), nil, msg, c)
}

// FailWithError 从验证错误中提取字段的自定义错误消息，obj为对应结构体
func FailWithError(err error, obj any, c *gin.Context) {
	msg := GetValidMsg(err, obj)
	FailWithMessage(domain.TypeUnknown, msg, c)
}

func GetValidMsg(err error, obj any) string {
	getObj := reflect.TypeOf(obj)
	var errs validator.ValidationErrors

	if errors.As(err, &errs) {
		for _, e := range errs {
			if f, exits := getObj.FieldByName(e.Field()); exits {
				//如果字段存在，通过 Tag.GetFirstPage("msg") 获取字段的自定义错误消息，并将其作为函数的返回值。
				return f.Tag.Get("msg")
			}
		}
	}

	return err.Error()
}
