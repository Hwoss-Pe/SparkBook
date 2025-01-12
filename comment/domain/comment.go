package domain

import "time"

type Comment struct {
	Id          int64  `json:"id"`
	Commentator User   `json:"user"` // 评论者
	Biz         string `json:"biz"`
	BizId       int64  `json:"BizId"`

	Content string `json:"content"`
	//根评论
	RootComment *Comment `json:"rootComment"`
	//父评论
	ParentComment *Comment  `json:"parentComment"`
	Children      []Comment `json:"children"`

	CTime time.Time `json:"ctime"`
	UTime time.Time `json:"utime"`
}

type User struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}
