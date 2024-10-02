package domain

type Article struct {
	ID      int64
	Title   string
	Content string
	Author  Author
	Status  StatusType
	Ctime   int64
	Utime   int64
}

type Author struct {
	Id   int64
	Name string
}
