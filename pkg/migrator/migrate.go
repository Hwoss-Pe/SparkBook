package migrator

//数据迁移的东西

type Entity interface {
	Id() int64
	CompareTo(dst Entity) bool
}
