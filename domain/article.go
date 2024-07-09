package domain

type Article struct {
	ID      int64
	Title   string
	Content string
	Author  Author
	Status  StatusType
}

type Author struct {
	Id   int64
	Name string
}
