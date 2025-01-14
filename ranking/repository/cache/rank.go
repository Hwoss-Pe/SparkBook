package cache

import (
	"Webook/ranking/domain"
	"encoding/json"
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

func (r *RedisRankingCache) Set(ctx context.Context, arts []domain.Article) error {
	//缓存的内容是缓存对应的摘要不是内容
	for _, art := range arts {
		art.Content = art.Abstract()
	}
	val, err := json.Marshal(arts)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, r.key, val, r.expiration).Err()
}

func (r *RedisRankingCache) Get(ctx context.Context) ([]domain.Article, error) {
	val, err := r.client.Get(ctx, r.key).Bytes()
	if err != nil {
		return nil, err
	}
	var res []domain.Article
	err = json.Unmarshal(val, &res)
	return res, err

}

func NewRedisRankingCache(client redis.Cmdable) *RedisRankingCache {
	return &RedisRankingCache{
		key:        "ranking:article",
		client:     client,
		expiration: time.Minute * 3,
	}
}
