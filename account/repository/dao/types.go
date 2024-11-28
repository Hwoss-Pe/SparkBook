package dao

import "context"

type AccountDAO interface {
	AddActivities(ctx context.Context, activities ...AccountActivity) error
}

//这里有一张账户本体表，还有个入账表

func (AccountActivity) TableName() string {
	return "account_activities"
}

type Account struct {
	Id  int64 `gorm:"primaryKey,autoIncrement" bson:"id,omitempty"`
	Uid int64 `gorm:"uniqueIndex:account_uid"`
	//一个uid和一个对外的账号
	Account int64 `gorm:"uniqueIndex:account_uid"`
	// 一个人可能有很多账号
	Type     uint8 `gorm:"uniqueIndex:account_uid"`
	Balance  int64
	Currency string
	Ctime    int64
	Utime    int64
}

type AccountActivity struct {
	Id    int64 `gorm:"primaryKey,autoIncrement" bson:"id,omitempty"`
	Uid   int64 `gorm:"index:account_uid"`
	Biz   string
	BizId int64
	// account 账号
	Account     int64 `gorm:"index:account_uid"`
	AccountType uint8 `gorm:"index:account_uid"`
	// 调整的金额，有些设计不想引入负数，就会增加一个类型
	// 标记是增加还是减少
	Amount   int64
	Currency string

	Ctime int64
	Utime int64
}
