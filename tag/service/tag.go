package service

import (
	"Webook/pkg/logger"
	"Webook/tag/domain"
	"Webook/tag/events"
	"Webook/tag/repository"
	"context"
)

type TagService interface {
	CreateTag(ctx context.Context, uid int64, name string) (int64, error)
	AttachTags(ctx context.Context, uid int64, biz string, bizId int64, tags []int64) error
	GetTags(ctx context.Context, uid int64) ([]domain.Tag, error)
	GetBizTags(ctx context.Context, uid int64, biz string, bizId int64) ([]domain.Tag, error)
}
type tagService struct {
	repo     repository.TagRepository
	logger   logger.Logger
	producer events.Producer
}

func (t *tagService) CreateTag(ctx context.Context, uid int64, name string) (int64, error) {
	return t.repo.CreateTag(ctx, domain.Tag{
		Uid:  uid,
		Name: name,
	})
}

func (t *tagService) AttachTags(ctx context.Context, uid int64, biz string, bizId int64, tags []int64) error {
	err := t.repo.BindTagToBiz(ctx, uid, biz, bizId, tags)
	if err != nil {
		return err
	}
	names, err := t.repo.GetTagsById(ctx, tags)
	if err != nil {
		return err
	}
	go func() {
		er := t.producer.ProduceSyncEvent(ctx, events.BizTags{
			Uid:   uid,
			Biz:   biz,
			BizId: bizId,
			Tags:  sliceNames(names),
		})
		if er != nil {
			t.logger.Error("发送标签事件失败",
				logger.Int64("biz_id", bizId),
				logger.Error(er))
		}
	}()
	return nil
}

func sliceNames(tags []domain.Tag) []string {
	res := make([]string, 0, len(tags))
	for _, t := range tags {
		res = append(res, t.Name)
	}
	return res
}

func (t *tagService) GetTags(ctx context.Context, uid int64) ([]domain.Tag, error) {
	return t.repo.GetTags(ctx, uid)
}

func (t *tagService) GetBizTags(ctx context.Context, uid int64, biz string, bizId int64) ([]domain.Tag, error) {
	return t.repo.GetBizTags(ctx, uid, biz, bizId)
}

func NewTagService(repo repository.TagRepository,
	producer events.Producer,
	l logger.Logger) TagService {
	return &tagService{
		producer: producer,
		repo:     repo,
		logger:   l,
	}
}
