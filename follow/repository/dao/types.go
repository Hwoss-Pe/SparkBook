package dao

import (
	"golang.org/x/net/context"
)

type FollowRelationDao interface {
	FollowRelationList(ctx context.Context, follower, offset, limit int64) ([]FollowRelation, error)
	FollowerRelationList(ctx context.Context, follower int64, followee int64) ([]FollowRelation, error)
	// CreateFollowRelation 创建联系人
	CreateFollowRelation(ctx context.Context, c FollowRelation) error
	// UpdateStatus 更新状态
	UpdateStatus(ctx context.Context, followee int64, follower int64, status uint8) error
	// CntFollower 统计计算关注自己的人有多少
	CntFollower(ctx context.Context, uid int64) (int64, error)
	// CntFollowee 统计自己关注了多少人
	CntFollowee(ctx context.Context, uid int64) (int64, error)
}

// UserRelation 用户关系
type UserRelation struct {
	ID     int64 `json:"id"`
	Uid1   int64 `json:"uid"`
	Uid2   int64 `json:"uid"`
	Block  bool  //拉黑
	Mute   bool  //屏蔽
	Follow bool  //关注
}
type FollowStatics struct {
	ID  int64 `gorm:"primaryKey,autoIncrement,column:id"`
	Uid int64 `gorm:"unique"`
	//关注/粉丝
	Followers int64
	Followees int64

	utime int64
	ctime int64
}

type FollowRelation struct {
	ID int64 `gorm:"primaryKey,autoIncrement,column:id"`

	Follower int64 `gorm:"type=int(11);not null; uniqueIndex:follow_followee"`
	Followee int64 `gorm:"type=int(11);not null; uniqueIndex:follow_followee"`

	Status uint8

	Ctime int64
	Utime int64
}
