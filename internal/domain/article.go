package domain

import "time"

type Article struct {
	ID      int64
	Title   string
	Content string
	Author  Author
	Status  StatusType
	Ctime   time.Time
	Utime   time.Time
}

type Author struct {
	Id   int64
	Name string
}
