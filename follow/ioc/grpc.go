package ioc

import (
	grpc2 "Webook/follow/grpc"
	"Webook/pkg/grpcx"
	"Webook/pkg/logger"
	"github.com/spf13/viper"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
)

func InitGRPCxServer(svc *grpc2.FollowServiceServer, ecli *clientv3.Client, logger logger.Logger) *grpcx.Server {
	type Config struct {
		Port     int    `yaml:"port"`
		EtcdAddr string `yaml:"etcdAddr"`
		EtcdTTL  int64  `yaml:"etcdTTL"`
	}
	var config Config
	err := viper.UnmarshalKey("grpc.server", &config)
	if err != nil {
		panic(err)
	}
	server := grpc.NewServer()
	svc.Register(server)
	return &grpcx.Server{
		Server:     server,
		Port:       config.Port,
		Name:       "follow",
		L:          logger,
		EtcdClient: ecli,
		EtcdTTL:    config.EtcdTTL,
	}
}
