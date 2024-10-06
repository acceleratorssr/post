package events

type Author struct {
	Id   uint64
	Name string
}

type Article struct {
	ID      uint64
	Title   string
	Content string
	Author  Author
	Ctime   int64
	Utime   int64
}

type PublishEvent struct {
	Article   *Article
	OnlyCache bool
	Uid       uint64
}
