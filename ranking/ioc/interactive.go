package ioc

import (
	intrv1 "Webook/api/proto/gen/api/proto/intr/v1"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func InitInterActiveRpcClient() intrv1.InteractiveServiceClient {
	type Config struct {
		Addr string `yaml:"addr"`
	}
	var config Config
	err := viper.UnmarshalKey("grpc.client.intr", &config)
	if err != nil {
		panic(err)
	}
	conn, err := grpc.Dial(config.Addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	client := intrv1.NewInteractiveServiceClient(conn)
	return client
}
