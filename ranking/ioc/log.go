package ioc

import (
	"Webook/pkg/logger"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func InitLogger() logger.Logger {
	config := zap.NewDevelopmentConfig()
	err := viper.UnmarshalKey("log", &config)
	if err != nil {
		panic(err)
	}

	build, err := config.Build()
	if err != nil {
		panic(err)
	}
	return logger.NewZapLogger(build)
}
