//go:build !k8s

package config

var Config = WebookConfig{
	DB: DBConfig{
		DSN: "root:root@tcp(localhost:13316)/webook",
	},
	Redis: RedisConfig{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	},
}
