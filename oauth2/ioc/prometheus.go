package ioc

import (
	"Webook/oauth2/service/wechat"
	"Webook/pkg/logger"
	"github.com/spf13/viper"
)

func InitPrometheus(logv1 logger.Logger) wechat.Service {
	svc := InitService(logv1)
	type Config struct {
		NameSpace  string `yaml:"nameSpace"`
		Subsystem  string `yaml:"subsystem"`
		InstanceID string `yaml:"instanceId"`
		Name       string `yaml:"name"`
	}
	var cfg Config
	err := viper.UnmarshalKey("prometheus", &cfg)
	if err != nil {
		panic(err)
	}
	return wechat.NewPrometheusDecorator(svc, cfg.NameSpace, cfg.Subsystem, cfg.InstanceID, cfg.Name)
}
