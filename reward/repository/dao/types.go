package dao

import "context"

type RewardDAO interface {
	Insert(ctx context.Context, r Reward) (int64, error)
	GetReward(ctx context.Context, rid int64) (Reward, error)
	UpdateStatus(ctx context.Context, rid int64, status uint8) error
}

type Reward struct {
	Id        int64  `gorm:"primaryKey,autoIncrement" bson:"id,omitempty"`
	biz       string `gorm:"index:biz_biz_id" bson:"biz"`
	BizId     int64  `gorm:"index:biz_biz_id"`
	BizName   string
	TargetUid int64 `gorm:"index"`
	// 直接采用 RewardStatus 的取值
	Status uint8

	// 打赏的人
	Uid    int64
	Amount int64
	Ctime  int64
	Utime  int64
}
