package ioc

import (
	"Webook/pkg/logger"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func InitLogger() logger.Logger {
	// 这里我们用一个小技巧，
	// 就是直接使用 zap 本身的配置结构体来处理
	cfg := zap.NewDevelopmentConfig()
	cfg.OutputPaths = []string{"/var/log/bff.log"}
	err := viper.UnmarshalKey("log", &cfg)
	if err != nil {
		panic(err)
	}
	l, err := cfg.Build()
	if err != nil {
		panic(err)
	}
	return logger.NewZapLogger(l)
}
