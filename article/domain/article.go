package domain

type Article struct {
	ID      uint64
	Title   string
	Content string
	Author  Author
	Status  StatusType
	Ctime   int64
	Utime   int64
}

type Author struct {
	Id   uint64
	Name string
}
