package cache

import (
	"Webook/user/domain"
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

var ErrKeyNotExist = redis.Nil

type UserCache interface {
	Delete(ctx context.Context, id int64) error
	Get(ctx context.Context, id int64) (domain.User, error)
	Set(ctx context.Context, u domain.User) error
}

type RedisUserCache struct {
	cmd        redis.Cmdable
	expiration time.Duration
}

func NewRedisUserCache(cmd redis.Cmdable) *RedisUserCache {
	return &RedisUserCache{cmd: cmd, expiration: time.Minute * 15}
}

func (r *RedisUserCache) Delete(ctx context.Context, id int64) error {
	return r.cmd.Del(ctx, r.Key(id)).Err()
}

func (r *RedisUserCache) Get(ctx context.Context, id int64) (domain.User, error) {
	key := r.Key(id)
	data, err := r.cmd.Get(ctx, key).Result()
	if err != nil {
		return domain.User{}, err
	}
	var du domain.User
	err = json.Unmarshal([]byte(data), &du)
	return du, err
}

func (r *RedisUserCache) Set(ctx context.Context, u domain.User) error {
	data, err := json.Marshal(u)
	if err != nil {
		return err
	}
	key := r.Key(u.Id)
	return r.cmd.Set(ctx, key, data, r.expiration).Err()
}

func (r *RedisUserCache) Key(id int64) string {
	return fmt.Sprintf("user:info:%d", id)
}
