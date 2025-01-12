package dao

import (
	"database/sql"
	"gorm.io/gorm"
)

type Comment struct {
	Id int64 `json:"id" gorm:"column:id;primaryKey"`
	//发表评论的用户
	Uid int64  `json:"uid" gorm:"column:uid;index"`
	Biz string `gorm:"column:biz;index:biz_type_id" json:"biz"`
	// 对应的业务ID
	BizID int64 `gorm:"column:biz_id;index:biz_type_id" json:"bizID"`

	//	根评论为0就是一级
	RootID sql.NullInt64 `json:"rootID" gorm:"column:root_id;index"`
	// 父级评论
	PID sql.NullInt64 `gorm:"column:pid;index" json:"pid"`
	//设置级联删除，删除的时候会把父评论对应的给一起删了,可以理解成PID找不到的id的时候就会把自己也删掉
	ParentComment *Comment `gorm:"ForeignKey:PID;AssociationForeignKey:ID;constraint:OnDelete:CASCADE"`
	// 评论内容
	Content string `gorm:"type:text;column:content" json:"content"`
	// 创建时间
	Ctime int64 `gorm:"column:ctime;" json:"ctime"`
	// 更新时间
	Utime int64 `gorm:"column:utime;" json:"utime"`
}

func (*Comment) TableName() string {
	return "comments"
}

type GORMCommentDAO struct {
	db *gorm.DB
}
