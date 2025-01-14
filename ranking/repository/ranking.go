package repository

import (
	"Webook/ranking/domain"
	"Webook/ranking/repository/cache"
	"github.com/ecodeclub/ekit/syncx/atomicx"
	"golang.org/x/net/context"
)

type RankingRepository interface {
	ReplaceTopN(ctx context.Context, arts []domain.Article) error
	GetTopN(ctx context.Context) ([]domain.Article, error)
}

func NewCachedRankingRepository(
	redisCache *cache.RedisRankingCache,
	localCache *cache.RankingLocalCache) RankingRepository {
	return &CachedRankingRepository{
		redisCache: redisCache,
		localCache: localCache,
	}
}

type CachedRankingRepository struct {
	redisCache *cache.RedisRankingCache
	localCache *cache.RankingLocalCache
	// 考虑将这个本地缓存塞进去 RankingCache 里面，作为一个实现
	topN atomicx.Value[[]domain.Article]
}

func (c *CachedRankingRepository) ReplaceTopN(ctx context.Context, arts []domain.Article) error {
	//先更不信本地缓存,再更新redis
	_ = c.localCache.Set(ctx, arts)
	return c.redisCache.Set(ctx, arts)
}

func (c *CachedRankingRepository) GetTopN(ctx context.Context) ([]domain.Article, error) {
	arts, err := c.localCache.Get(ctx)
	if err == nil {
		return arts, nil
	}
	// 回写本地缓存
	arts, err = c.redisCache.Get(ctx)
	if err == nil {
		_ = c.localCache.Set(ctx, arts)
	} else {
		// 这里，我们没有进一步区分是什么原因导致的 Redis 错误
		return c.localCache.ForceGet(ctx)
	}
	return arts, err
}
