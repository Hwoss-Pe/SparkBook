package repository

import (
	"Webook/interactive/repository/dao"
	"context"
)

type NotificationRepository interface {
	List(ctx context.Context, uid int64, typ string, offset int, limit int) ([]dao.Notification, error)
	MarkReadByIds(ctx context.Context, uid int64, ids []int64) error
	MarkReadByType(ctx context.Context, uid int64, typ string) error
	CountUnreadAll(ctx context.Context, uid int64) (interaction int64, follow int64, system int64, total int64, err error)
}

type GORMNotificationRepository struct {
	dao dao.NotificationDAO
}

func NewNotificationRepository(d dao.NotificationDAO) *GORMNotificationRepository {
	return &GORMNotificationRepository{dao: d}
}

func (r *GORMNotificationRepository) List(ctx context.Context, uid int64, typ string, offset int, limit int) ([]dao.Notification, error) {
	return r.dao.ListByTypeAndUser(ctx, uid, typ, offset, limit)
}

func (r *GORMNotificationRepository) MarkReadByIds(ctx context.Context, uid int64, ids []int64) error {
	return r.dao.MarkReadByIds(ctx, uid, ids)
}

func (r *GORMNotificationRepository) MarkReadByType(ctx context.Context, uid int64, typ string) error {
	return r.dao.MarkReadByType(ctx, uid, typ)
}

func (r *GORMNotificationRepository) CountUnreadAll(ctx context.Context, uid int64) (int64, int64, int64, int64, error) {
	return r.dao.CountUnreadAll(ctx, uid)
}
