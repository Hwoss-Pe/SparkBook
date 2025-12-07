package domain

import "time"

type Article struct {
	Id         int64
	Title      string
	Status     ArticleStatus
	Content    string
	CoverImage string // 封面图片URL

	Author Author
	Ctime  time.Time
	Utime  time.Time
}

const (
	ArticleStatusUnknown = iota
	ArticleStatusUnpublished
	ArticleStatusPublished
	ArticleStatusPrivate
)

func (a Article) Published() bool {
	return a.Status == ArticleStatusPublished
}

type ArticleStatus uint8

func (s ArticleStatus) ToUint8() uint8 {
	return uint8(s)
}

func (a Article) Abstract() string {
	runes := []rune(a.Content)
	if len(runes) < 100 {
		return a.Content
	} else {
		return string(runes[:100])
	}
}

type Author struct {
	Id   int64
	Name string
}
