package domain

type Article struct {
	ID      uint64 // 唯一标识
	Title   string
	Content string
	Author  Author
	Ctime   int64
	Utime   int64
}

type Author struct {
	Id   uint64
	Name string
}

type List struct {
	Limit     int
	LastValue int64 // 保存在客户端，用于翻页时防重复数据
	Desc      bool  // 0为降序，1为升序
	OrderBy   string
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
