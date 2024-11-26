package ioc

import (
	smsv1 "Webook/api/proto/gen/api/proto/sms/v1"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func InitSmsRpcClient() smsv1.SmsServiceClient {
	type config struct {
		Target string `yaml:"target"`
	}
	var cfg config
	err := viper.UnmarshalKey("grpc.client.sms", &cfg)
	if err != nil {
		panic(err)
	}
	conn, err := grpc.Dial(cfg.Target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	return smsv1.NewSmsServiceClient(conn)
}
