package cache

import (
	"Webook/follow/domain"
	"fmt"
	"github.com/redis/go-redis/v9"
	"golang.org/x/net/context"
	"strconv"
)

const (
	// 被多少人关注
	fieldFollowerCnt = "follower_cnt"
	// 关注了多少人
	fieldFolloweeCnt = "followee_cnt"
)

type FollowCache interface {
	StaticsInfo(ctx context.Context, uid int64) (domain.FollowStatics, error)
	SetStaticsInfo(ctx context.Context, uid int64, statics domain.FollowStatics) error
	Follow(ctx context.Context, follow, followee int64) error
	CancelFollow(ctx context.Context, follower, followee int64) error
}

var ErrKeyNotExist = redis.Nil

type RedisFollowCache struct {
	client redis.Cmdable
}

func (r *RedisFollowCache) StaticsInfo(ctx context.Context, uid int64) (domain.FollowStatics, error) {
	data, err := r.client.HGetAll(ctx, r.staticsKey(uid)).Result()
	if err != nil {
		return domain.FollowStatics{}, err
	}
	if len(data) == 0 {
		return domain.FollowStatics{}, ErrKeyNotExist
	}
	followerCnt, _ := strconv.ParseInt(data[fieldFollowerCnt], 10, 64)
	followeeCnt, _ := strconv.ParseInt(data[fieldFollowerCnt], 10, 64)
	return domain.FollowStatics{
		Followees: followeeCnt,
		Followers: followerCnt,
	}, nil
}

func (r *RedisFollowCache) SetStaticsInfo(ctx context.Context, uid int64, statics domain.FollowStatics) error {
	key := r.staticsKey(uid)
	return r.client.HMSet(ctx, key, fieldFolloweeCnt, statics.Followees, fieldFollowerCnt, statics.Followers).Err()
}

func (r *RedisFollowCache) Follow(ctx context.Context, follow, followee int64) error {
	return r.updateStaticsInfo(ctx, follow, followee, 1)
}

func (r *RedisFollowCache) CancelFollow(ctx context.Context, follower, followee int64) error {
	return r.updateStaticsInfo(ctx, follower, followee, -1)
}

func (r *RedisFollowCache) updateStaticsInfo(ctx context.Context, follower, followee int64, delta int64) error {
	//	这里对于redis也是需要并发的操作，关注就相等于给人加粉丝，并且自己的关注数量加1
	tx := r.client.TxPipeline()
	// 增加 follower 的关注多少人的数量
	tx.HIncrBy(ctx, r.staticsKey(follower), fieldFolloweeCnt, delta)
	// 增加 followee 被多少人关注的数量
	tx.HIncrBy(ctx, r.staticsKey(followee), fieldFollowerCnt, delta)
	_, err := tx.Exec(ctx)
	return err
}

func (r *RedisFollowCache) staticsKey(uid int64) string {
	return fmt.Sprintf("follow:statics:%d", uid)
}

func NewRedisFollowCache(client redis.Cmdable) FollowCache {
	return &RedisFollowCache{
		client: client,
	}
}
