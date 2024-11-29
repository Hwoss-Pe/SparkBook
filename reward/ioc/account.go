package ioc

import (
	accountv1 "Webook/api/proto/gen/api/proto/account/v1"
	"github.com/spf13/viper"
	etcdv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/resolver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func InitAccountClient(etcdClient *etcdv3.Client) accountv1.AccountServiceClient {
	type Config struct {
		Target string `json:"target"`
		Secure bool   `json:"secure"`
	}
	var cfg Config
	err := viper.UnmarshalKey("grpc.client.account", &cfg)
	if err != nil {
		panic(err)
	}
	//创建etcd解析器然后给grpc配置链接
	builder, err := resolver.NewBuilder(etcdClient)
	if err != nil {
		panic(err)
	}
	options := []grpc.DialOption{grpc.WithResolvers(builder)}
	if !cfg.Secure {
		options = append(options, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}
	cc, err := grpc.Dial(cfg.Target, options...)
	if err != nil {
		panic(err)
	}
	return accountv1.NewAccountServiceClient(cc)
}
