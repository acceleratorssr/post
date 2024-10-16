package events

type Author struct {
	Id   uint64 `json:"id"`
	Name string `json:"name"`
}

type Article struct {
	ID      uint64 `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Author  Author `json:"author"`
	Ctime   int64  `json:"ctime"`
	Utime   int64  `json:"utime"`
}

type PublishEvent struct {
	Article   *Article `json:"article"`
	OnlyCache bool     `json:"only_cache"`
	Uid       uint64   `json:"uid"`
	Delete    bool     `json:"delete"`
}

type RecommendEvent struct {
	FeedbackType string `json:"feedback_type"`
	ArticleID    uint64 `json:"article_id"`
	UserId       string `json:"uid"`
	ItemId       string `json:"item_id"`
	Timestamp    string `json:"timestamp"`
}
