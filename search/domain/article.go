package domain

type Article struct {
	Id      int64
	Title   string
	Status  int32
	Content string
	Tags    []string
}
