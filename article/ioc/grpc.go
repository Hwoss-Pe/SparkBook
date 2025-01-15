package ioc

import (
	grpc2 "Webook/article/grpc"
	"Webook/pkg/grpcx"
	"Webook/pkg/logger"
	"github.com/spf13/viper"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
)

func InitGRPCxServer(articleServer *grpc2.ArticleServiceServer,
	l logger.Logger, etcd *clientv3.Client) *grpcx.Server {
	type Config struct {
		Port     int    `yaml:"port"`
		EtcdAddr string `yaml:"etcdAddr"`
		EtcdTTL  int64  `yaml:"etcdTTL"`
	}
	var cfg Config
	err := viper.UnmarshalKey("grpc.server", &cfg)
	if err != nil {
		panic(err)
	}
	server := grpc.NewServer()
	articleServer.Register(server)
	return &grpcx.Server{
		Server:     server,
		Port:       cfg.Port,
		Name:       "article",
		L:          l,
		EtcdClient: etcd,
		EtcdTTL:    cfg.EtcdTTL,
	}
}
