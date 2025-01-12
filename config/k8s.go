//go:build k8s

package config

var Config = WebookConfig{
	DB: DBConfig{
		DSN: "root:root@tcp(webook-mysql:3308)/webook",
	},
	Redis: RedisConfig{
		Addr:     "webook-redis:6379",
		Password: "",
		DB:       1,
	},
}
