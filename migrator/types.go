package migrator

type Entity interface {
	GetID() int64
	CompareWith(e Entity) bool
}

//type CompareWith interface {
//	CompareWith(e Entity) bool
//}
