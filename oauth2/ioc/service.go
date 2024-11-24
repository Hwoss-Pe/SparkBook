package ioc

import (
	"Webook/oauth2/service/wechat"
	"Webook/pkg/logger"
	"github.com/spf13/viper"
)

func InitService(log logger.Logger) wechat.Service {
	type Config struct {
		AppID     string `yaml:"appId"`
		AppSecret string `yaml:"appSecret"`
	}
	var cfg Config
	err := viper.UnmarshalKey("weChatConf", &cfg)
	if err != nil {
		panic(err)
	}
	return wechat.NewService(cfg.AppID, cfg.AppSecret, log)
}
