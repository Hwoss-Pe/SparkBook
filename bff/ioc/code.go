package ioc

import (
	codev1 "Webook/api/proto/gen/api/proto/code/v1"
	"github.com/spf13/viper"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/resolver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func InitCodeClient(ecli *clientv3.Client) codev1.CodeServiceClient {
	type Config struct {
		Target string `json:"target"`
		Secure bool   `json:"secure"`
	}
	var cfg Config
	err := viper.UnmarshalKey("grpc.client.code", &cfg)
	if err != nil {
		panic(err)
	}
	rs, err := resolver.NewBuilder(ecli)
	if err != nil {
		panic(err)
	}
	opts := []grpc.DialOption{grpc.WithResolvers(rs)}
	if !cfg.Secure {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}
	cc, err := grpc.Dial(cfg.Target, opts...)
	if err != nil {
		panic(err)
	}
	return codev1.NewCodeServiceClient(cc)
}
