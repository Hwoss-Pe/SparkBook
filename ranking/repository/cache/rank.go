package cache

import (
	"Webook/ranking/domain"
	"github.com/redis/go-redis/v9"
	"golang.org/x/net/context"
	"time"
)

type RankingCache interface {
	Set(ctx context.Context, arts []domain.Article) error
	Get(ctx context.Context) ([]domain.Article, error)
}

type RedisRankingCache struct {
	client     redis.Cmdable
	key        string
	expiration time.Duration
}
