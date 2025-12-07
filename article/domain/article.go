package domain

import "time"

type Article struct {
	Id         int64
	Title      string
	Content    string
	CoverImage string // 封面图片URL
	Status     ArticleStatus
	Ctime      time.Time
	Utime      time.Time

	Author Author
}

func (a Article) Abstract() string {
	//这里截取字节可能会出现中文或者表情包的问题，需要用rune进行一个unicode的字节统计
	runes := []rune(a.Content)
	if len(runes) < 100 {
		return a.Content
	}
	return string(runes[:100])
}

type ArticleStatus uint8

const (
	//	未发表，已发表，仅自己可见

	ArticleStatusUnknown ArticleStatus = iota
	ArticleStatusUnpublished
	ArticleStatusPublished
	ArticleStatusPrivate
)

// Published 出版更新
func (a Article) Published() bool {
	return a.Status == ArticleStatusPublished
}
func (s ArticleStatus) ToUint8() uint8 {
	return uint8(s)
}

// Author 在领域的概念里面，他只有作者而不是用户的概念，用户有额外的用户领域
type Author struct {
	Id     int64
	Name   string
	Avatar string // 作者头像URL
}
