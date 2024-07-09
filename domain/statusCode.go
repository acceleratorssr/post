package domain

type StatusType uint16

const (
	TypeUnknown StatusType = iota // 注意，一般情况下避免零值问题，0对应未知情况
	TypeSaved
	TypePublished
	TypeWithdraw
	ErrPostArticle
	ErrSystem StatusType = 555
)

var (
	ErrorMap = map[StatusType]string{
		TypeUnknown:    "TypeUnknown",
		TypeSaved:      "文章已保存",
		TypePublished:  "文章已发表",
		ErrPostArticle: "发布文章失败",
	}
)

func (e StatusType) string() string {
	switch e {
	case TypeSaved:
		return "文章已保存"
	case TypePublished:
		return "文章已发表"
	case ErrPostArticle:
		return "发布文章失败"
	default:
		return "TypeUnknown"
	}
}

func (e StatusType) ToUint8() uint8 {
	t := e
	return uint8(t)
}

func (e StatusType) ToInt() int {
	t := e
	return int(t)
}
