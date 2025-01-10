package service

import (
	"Webook/interactive/domain"
	"Webook/interactive/repository"
	"Webook/pkg/logger"
	"context"
	"golang.org/x/sync/errgroup"
)

type InteractiveService interface {
	IncrReadCnt(ctx context.Context, biz string, bizId int64) error
	// Like 点赞
	Like(ctx context.Context, biz string, bizId int64, uid int64) error
	// CancelLike 取消点赞
	CancelLike(ctx context.Context, biz string, bizId int64, uid int64) error
	// Collect 收藏
	Collect(ctx context.Context, biz string, bizId, cid, uid int64) error
	Get(ctx context.Context, biz string, bizId, uid int64) (domain.Interactive, error)
	GetByIds(ctx context.Context, biz string, bizIds []int64) (map[int64]domain.Interactive, error)
}

type interactiveService struct {
	repo repository.InteractiveRepository
	l    logger.Logger
}

func NewInteractiveService(repo repository.InteractiveRepository,
	l logger.Logger) InteractiveService {
	return &interactiveService{
		repo: repo,
		l:    l,
	}
}

func (i *interactiveService) IncrReadCnt(ctx context.Context, biz string, bizId int64) error {
	return i.repo.IncrReadCnt(ctx, biz, bizId)
}

func (i *interactiveService) Like(ctx context.Context, biz string, bizId int64, uid int64) error {
	return i.repo.IncrLike(ctx, biz, bizId, uid)
}

func (i *interactiveService) CancelLike(ctx context.Context, biz string, bizId int64, uid int64) error {
	return i.repo.DecrLike(ctx, biz, bizId, uid)
}

func (i *interactiveService) Collect(ctx context.Context, biz string, bizId, cid, uid int64) error {
	return i.repo.AddCollectionItem(ctx, biz, bizId, cid, uid)
}

func (i *interactiveService) Get(ctx context.Context, biz string, bizId, uid int64) (domain.Interactive, error) {
	intr, err := i.repo.Get(ctx, biz, bizId)
	if err != nil {
		return domain.Interactive{}, err
	}
	var eg errgroup.Group
	eg.Go(func() error {
		intr.Liked, err = i.repo.Liked(ctx, biz, bizId, uid)
		return err
	})
	eg.Go(func() error {
		intr.Liked, err = i.repo.Collected(ctx, biz, bizId, uid)
		return err
	})
	err = eg.Wait()
	if err != nil {
		// 这个查询失败只需要记录日志就可以，不需要中断执行
		i.l.Error("查询用户是否点赞的信息失败",
			logger.String("biz", biz),
			logger.Int64("bizId", bizId),
			logger.Int64("uid", uid),
			logger.Error(err))
	}
	return intr, nil
}

func (i *interactiveService) GetByIds(ctx context.Context, biz string, bizIds []int64) (map[int64]domain.Interactive, error) {
	intrs, err := i.repo.GetByIds(ctx, biz, bizIds)
	if err != nil {
		return nil, err
	}
	//返回的是一个map，并且用bizId作为键
	res := make(map[int64]domain.Interactive, len(intrs))
	for _, intr := range intrs {
		res[intr.BizId] = intr
	}
	return res, nil
}
