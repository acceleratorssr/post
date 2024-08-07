package utils

// （模块+错误方+错误编号） xx + x + xxx
// InvalidInput = 014001//前导0为8进制

// TODO 和statusCode好像有点重复，考虑合并
const (
	UserInvalidInput  = 104001
	UserOrPasswordErr = 104002

	UserServiceErr = 105001
)

const (
	ArticleInvalidInput = 114001

	ArticleServiceErr = 115001
)
