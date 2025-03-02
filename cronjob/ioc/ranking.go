package ioc

import (
	rankingv1 "Webook/api/proto/gen/api/proto/ranking/v1"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func InitRankingRpcClient() rankingv1.RankingServiceClient {
	type Config struct {
		Addr string `yaml:"addr"`
	}
	var config Config
	err := viper.UnmarshalKey("grpc.client.ranking", &config)
	if err != nil {
		panic(err)
	}
	conn, err := grpc.Dial(config.Addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	client := rankingv1.NewRankingServiceClient(conn)
	return client
}
