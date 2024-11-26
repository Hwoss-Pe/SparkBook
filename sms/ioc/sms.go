package ioc

import (
	"Webook/sms/service"
	"Webook/sms/service/localsms"
	"Webook/sms/service/tencent"
	"github.com/spf13/viper"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	tencentSMS "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
)

func InitSmsTencentService() service.Service {
	type Config struct {
		SecretID  string `yaml:"secretId"`
		SecretKey string `yaml:"secretKey"`
	}
	var cfg Config
	err := viper.UnmarshalKey("tencentSms", &cfg)
	c, err := tencentSMS.NewClient(common.NewCredential(cfg.SecretID, cfg.SecretKey),
		"ap-nanjing",
		profile.NewClientProfile())
	if err != nil {
		panic(err)
	}
	return tencent.NewService(c, "1400842696", "XXXXX")
}

func InitSmsService() service.Service {
	//return initSmsTencentService()
	return InitSmsMemoryService()
}

// InitSmsMemoryService 使用基于内存，输出到控制台的实现
func InitSmsMemoryService() service.Service {
	return localsms.NewService()
}
