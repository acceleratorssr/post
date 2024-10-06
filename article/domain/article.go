package domain

import "strconv"

type Article struct {
	ID      uint64 `json:"id"` // 唯一标识
	Title   string `json:"title"`
	Content string `json:"content"`
	Author  Author `json:"author"`
	Ctime   int64  `json:"ctime"`
	Utime   int64  `json:"utime"`
}

type Author struct {
	Id   uint64 `json:"id"`
	Name string `json:"name"`
}

type List struct {
	Limit     int    `json:"Limit"`
	LastValue int64  `json:"last_value"` // 保存在客户端，用于翻页时防重复数据
	Desc      bool   `json:"desc"`       // 0为升序，1为降序
	OrderBy   string `json:"order_by"`
}

func (l *List) Sort(limit int, orderBy string, desc bool, a []Article, b []Article) []Article {
	compare := func(x, y Article) bool {
		var xv, yv int64

		switch orderBy {
		case "Ctime":
			xv, yv = x.Ctime, y.Utime
		case "Utime":
			xv, yv = x.Utime, y.Utime
		default:
			return false
		}

		if desc {
			return xv < yv
		}
		return xv > yv
	}

	result := make([]Article, 0, limit)
	i, j := 0, 0
	for i < len(a) && j < len(b) && len(result) < limit {
		if compare(a[i], b[j]) {
			result = append(result, a[i])
			i++
		} else {
			result = append(result, b[j])
			j++
		}
	}
	for ; i < len(a) && len(result) < limit; i++ {
		result = append(result, a[i])
	}
	for ; j < len(b) && len(result) < limit; j++ {
		result = append(result, b[j])
	}

	return result
}

func GetArtCacheKey(aid uint64) string {
	return "article:" + strconv.FormatUint(aid, 10)
}

func KeyArtTopNBrief() string {
	return "article_rank_brief"
}
