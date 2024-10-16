package events

type ReadEvent struct {
	ID  uint64
	Uid uint64
	Aid uint64
}

type RecommendEvent struct {
	FeedbackType string `json:"feedback_type"`
	ArticleID    uint64 `json:"article_id"`
	UserId       string `json:"uid"`
	ItemId       string `json:"item_id"`
	Timestamp    string `json:"timestamp"`
}
