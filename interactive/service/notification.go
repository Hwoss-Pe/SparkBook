package service

import (
	intrv1 "Webook/api/proto/gen/api/proto/intr/v1"
	"Webook/interactive/repository"
	"context"
	"time"
)

type NotificationService interface {
	GetUnreadCounts(ctx context.Context, uid int64) (*intrv1.GetUnreadCountsResponse, error)
	GetNotifications(ctx context.Context, uid int64, typ string, offset int, limit int) (*intrv1.GetNotificationsResponse, error)
	MarkRead(ctx context.Context, uid int64, ids []int64, typ string) error
}

type notificationService struct {
	repo repository.NotificationRepository
}

func NewNotificationService(repo repository.NotificationRepository) NotificationService {
	return &notificationService{repo: repo}
}

func (s *notificationService) GetUnreadCounts(ctx context.Context, uid int64) (*intrv1.GetUnreadCountsResponse, error) {
	i, f, sy, total, err := s.repo.CountUnreadAll(ctx, uid)
	if err != nil {
		return nil, err
	}
	return &intrv1.GetUnreadCountsResponse{Interaction: int32(i), Follow: int32(f), System: int32(sy), Total: int32(total)}, nil
}

func (s *notificationService) GetNotifications(ctx context.Context, uid int64, typ string, offset int, limit int) (*intrv1.GetNotificationsResponse, error) {
	list, err := s.repo.List(ctx, uid, typ, offset, limit)
	if err != nil {
		return nil, err
	}
	items := make([]*intrv1.NotificationItem, 0, len(list))
	for _, n := range list {
		statusStr := "unread"
		if n.Status == 1 {
			statusStr = "read"
		}
		it := &intrv1.NotificationItem{
			Id:       n.Id,
			Category: n.Type,
			Content:  n.Content,
			Time:     time.Unix(n.Ctime, 0).Format(time.RFC3339),
			Status:   statusStr,
		}
		// 保证前端不空指针：填充 sender
		it.Sender = &intrv1.NotificationSender{Id: n.SenderId}
		// 根据业务类型填充 target 基本信息
		if n.BizType != "" && n.BizId > 0 {
			it.Target = &intrv1.NotificationTarget{Type: n.BizType, Id: n.BizId}
		}
		items = append(items, it)
	}
	return &intrv1.GetNotificationsResponse{Items: items}, nil
}

func (s *notificationService) MarkRead(ctx context.Context, uid int64, ids []int64, typ string) error {
	if len(ids) > 0 {
		return s.repo.MarkReadByIds(ctx, uid, ids)
	}
	if typ != "" {
		return s.repo.MarkReadByType(ctx, uid, typ)
	}
	return nil
}
