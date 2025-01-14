package cache

import (
	"Webook/ranking/domain"
	"github.com/ecodeclub/ekit/syncx/atomicx"
	"time"
)

// RankingLocalCache 因为本身数据只有一份，所以不需要借助真正的本地缓存
type RankingLocalCache struct {
	topN       *atomicx.Value[[]domain.Article]
	ddl        *atomicx.Value[time.Time]
	expiration time.Duration
}
