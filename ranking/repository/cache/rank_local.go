package cache

import (
	"Webook/ranking/domain"
	"errors"
	"github.com/ecodeclub/ekit/syncx/atomicx"
	"golang.org/x/net/context"
	"time"
)

// RankingLocalCache 因为本身数据只有一份，所以不需要借助真正的本地缓存
type RankingLocalCache struct {
	topN       *atomicx.Value[[]domain.Article]
	ddl        *atomicx.Value[time.Time]
	expiration time.Duration
}

func NewRankingLocalCache() *RankingLocalCache {
	return &RankingLocalCache{
		topN:       atomicx.NewValue[[]domain.Article](),
		ddl:        atomicx.NewValueOf[time.Time](time.Now()),
		expiration: time.Minute * 3,
	}
}
func (r *RankingLocalCache) Set(_ context.Context, arts []domain.Article) error {
	r.ddl.Store(time.Now().Add(time.Minute * 3))
	r.topN.Store(arts)
	return nil
}
func (r *RankingLocalCache) ForceGet(_ context.Context) ([]domain.Article, error) {
	return r.topN.Load(), nil
}
func (r *RankingLocalCache) Get(_ context.Context) ([]domain.Article, error) {
	arts := r.topN.Load()
	//如果本地缓存得到
	if len(arts) == 0 || r.ddl.Load().Before(time.Now()) {
		return nil, errors.New("本地缓存失效了")
	}
	return arts, nil
}
