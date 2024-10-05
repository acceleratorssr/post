package gin_ex

import (
	"fmt"
	"strconv"
)

type Code uint32

const (
	// Unknown 未知错误
	OK Code = 510

	// Canceled 用户取消操作
	Canceled Code = 511

	Unknown Code = 512

	// InvalidArgument 参数无效
	InvalidArgument Code = 513

	// DeadlineExceeded 超时
	DeadlineExceeded Code = 514

	// NotFound 资源未找到
	NotFound Code = 515

	// AlreadyExists 唯一实体已存在
	AlreadyExists Code = 516

	// PermissionDenied 通过认证但权限不足
	PermissionDenied Code = 517

	// ResourceExhausted 资源耗尽，也可能是触发降级等
	ResourceExhausted Code = 518

	// FailedPrecondition 操作前提条件不满足，应修改请求，不应重试
	FailedPrecondition Code = 519

	// Aborted 出现并发等冲突导致操作中止
	Aborted Code = 520

	// Internal 内部错误，非代码问题
	Internal Code = 521

	// Unavailable 不可用表示服务当前不可用
	// 这很可能是一种暂时的情况，可以通过回退重试来纠正
	// 需要注意，非幂等操作不安全
	Unavailable Code = 522

	// Unauthenticated 未通过认证
	Unauthenticated Code = 523

	// System 返回给客户端，无脑的错误
	System Code = 524

	_maxCode = 525
)

var strToCode = map[string]Code{
	"OK":                 OK,
	"Canceled":           Canceled,
	"Unknown":            Unknown,
	"InvalidArgument":    InvalidArgument,
	"DeadlineExceeded":   DeadlineExceeded,
	"NotFound":           NotFound,
	"AlreadyExists":      AlreadyExists,
	"PermissionDenied":   PermissionDenied,
	"ResourceExhausted":  ResourceExhausted,
	"FailedPrecondition": FailedPrecondition,
	"Aborted":            Aborted,
	"Internal":           Internal,
	"Unavailable":        Unavailable,
	"Unauthenticated":    Unauthenticated,
	"System":             System,
}

// UnmarshalJSON 当 JSON 数据中的某个字段需要解码为 Code 类型时，
// json.Unmarshal 函数会调用该自定义方法
func (c *Code) UnmarshalJSON(b []byte) error {
	if string(b) == "null" {
		return nil
	}

	// 防止空指针
	if c == nil {
		return fmt.Errorf("nil receiver passed to UnmarshalJSON")
	}

	stripped := string(b)
	if stripped[0] == '"' && stripped[len(stripped)-1] == '"' {
		stripped = stripped[1 : len(stripped)-1]
	}

	if ci, err := strconv.ParseUint(stripped, 10, 32); err == nil {
		if ci >= _maxCode {
			return fmt.Errorf("invalid code: %d", ci)
		}

		*c = Code(ci)
		return nil
	}

	if jc, ok := strToCode[stripped]; ok {
		*c = jc
		return nil
	}
	return fmt.Errorf("invalid code: %q", stripped)
}

func (c Code) ToInt() int {
	return int(c)
}
