package repository

import (
	"Webook/reward/domain"
	"context"
)

type RewardRepository interface {
	CreateReward(ctx context.Context, reward domain.Reward) (int64, error)
	GetReward(ctx context.Context, rid int64) (domain.Reward, error)

	// GetCachedCodeURL 这两个方法的名字明确带上了缓存的字眼
	// 是希望调用者明白这个是缓存下来的，属于业务逻辑的一部分
	GetCachedCodeURL(ctx context.Context, r domain.Reward) (domain.CodeURL, error)
	CachedCodeURL(ctx context.Context, cu domain.CodeURL, r domain.Reward) error
	UpdateStatus(ctx context.Context, rid int64, status domain.RewardStatus) error
}
