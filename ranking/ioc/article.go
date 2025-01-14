package ioc

import (
	articlev1 "Webook/api/proto/gen/api/proto/article/v1"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func InitArticleRpcClient() articlev1.ArticleServiceClient {
	type Config struct {
		Addr string `yaml:"addr"`
	}
	var config Config
	err := viper.UnmarshalKey("grpc.client.article", &config)
	if err != nil {
		panic(err)
	}
	conn, err := grpc.Dial(config.Addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	client := articlev1.NewArticleServiceClient(conn)
	return client
}
