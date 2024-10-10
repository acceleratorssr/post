package domain

import "strconv"

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

func KeyIncrReadCount(ObjType string, ObjID uint64) string {
	return "incr_read_count:" + ObjType + ":" + strconv.FormatUint(ObjID, 10)
}

func KeyIncrLikeCount(objType string, objID uint64) string {
	return "incr_Like_count:" + objType + ":" + strconv.FormatUint(objID, 10)
}
