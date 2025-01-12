package ioc

import (
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

func InitRedisClient() redis.Cmdable {
	// 实现 Redis 客户端的初始化
	client := redis.NewClient(&redis.Options{
		Addr: viper.GetString("redis.addr"),
	})
	return client
}
