package domain

type Like struct {
	ID        uint64
	LikeCount int64
	Ctime     int64
}

type List struct {
	Limit     int    `json:"Limit"`
	LastValue int64  `json:"last_value"` // 保存在客户端，用于翻页时防重复数据
	Desc      bool   `json:"desc"`       // 0为升序，1为降序
	OrderBy   string `json:"order_by"`
}
