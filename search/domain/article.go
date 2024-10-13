package domain

type Author struct {
	Id   uint64 `json:"id"`
	Name string `json:"name"`
}

type Article struct {
	ID      uint64   `json:"id"`
	Title   string   `json:"title"`
	Content string   `json:"content"`
	Author  Author   `json:"author"`
	Tags    []string `json:"tags"`
	Ctime   int64    `json:"ctime"`
	Utime   int64    `json:"utime"`
}

type ArticleEvent struct {
	Article   *Article `json:"article"`
	OnlyCache bool     `json:"only_cache"`
	Uid       uint64   `json:"uid"`
}
