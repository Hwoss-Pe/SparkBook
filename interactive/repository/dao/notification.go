package dao

import (
	"context"

	"gorm.io/gorm"
)

type NotificationDAO interface {
	Insert(ctx context.Context, n Notification) error
	ListByTypeAndUser(ctx context.Context, uid int64, typ string, offset int, limit int) ([]Notification, error)
	MarkReadByIds(ctx context.Context, uid int64, ids []int64) error
	MarkReadByType(ctx context.Context, uid int64, typ string) error
	CountUnreadByType(ctx context.Context, uid int64, typ string) (int64, error)
	CountUnreadAll(ctx context.Context, uid int64) (interaction int64, follow int64, system int64, total int64, err error)
}

type GORMNotificationDAO struct {
	db *gorm.DB
}

func NewGORMNotificationDAO(db *gorm.DB) *GORMNotificationDAO {
	return &GORMNotificationDAO{db: db}
}

type Notification struct {
	Id         int64  `gorm:"primaryKey,autoIncrement"`
	ReceiverId int64  `gorm:"index"`
	SenderId   int64  `gorm:"index"`
	Type       string `gorm:"type:varchar(64);index"`
	BizType    string `gorm:"type:varchar(64)"`
	BizId      int64  `gorm:"index"`
	Content    string `gorm:"type:text"`
	Status     uint8  `gorm:"type:tinyint;index"` // 0 未读, 1 已读
	Ctime      int64
	Utime      int64
}

func (d *GORMNotificationDAO) Insert(ctx context.Context, n Notification) error {
	return d.db.WithContext(ctx).Create(&n).Error
}

func (d *GORMNotificationDAO) ListByTypeAndUser(ctx context.Context, uid int64, typ string, offset int, limit int) ([]Notification, error) {
	var res []Notification
	err := d.db.WithContext(ctx).
		Where("receiver_id = ? AND type = ?", uid, typ).
		Order("ctime DESC").
		Offset(offset).Limit(limit).
		Find(&res).Error
	return res, err
}

func (d *GORMNotificationDAO) MarkReadByIds(ctx context.Context, uid int64, ids []int64) error {
	if len(ids) == 0 {
		return nil
	}
	return d.db.WithContext(ctx).
		Model(&Notification{}).
		Where("receiver_id = ? AND id IN ? AND status = ?", uid, ids, 0).
		Updates(map[string]any{"status": 1}).Error
}

func (d *GORMNotificationDAO) MarkReadByType(ctx context.Context, uid int64, typ string) error {
	return d.db.WithContext(ctx).
		Model(&Notification{}).
		Where("receiver_id = ? AND type = ? AND status = ?", uid, typ, 0).
		Updates(map[string]any{"status": 1}).Error
}

func (d *GORMNotificationDAO) CountUnreadByType(ctx context.Context, uid int64, typ string) (int64, error) {
	var cnt int64
	err := d.db.WithContext(ctx).
		Model(&Notification{}).
		Where("receiver_id = ? AND type = ? AND status = ?", uid, typ, 0).
		Count(&cnt).Error
	return cnt, err
}

func (d *GORMNotificationDAO) CountUnreadAll(ctx context.Context, uid int64) (int64, int64, int64, int64, error) {
	types := []string{"interaction", "follow", "system"}
	var vals [3]int64
	for i, t := range types {
		var cnt int64
		err := d.db.WithContext(ctx).
			Model(&Notification{}).
			Where("receiver_id = ? AND type = ? AND status = ?", uid, t, 0).
			Count(&cnt).Error
		if err != nil {
			return 0, 0, 0, 0, err
		}
		vals[i] = cnt
	}
	total := vals[0] + vals[1] + vals[2]
	return vals[0], vals[1], vals[2], total, nil
}
