package ioc

import (
	"Webook/pkg/grpcx"
	"Webook/pkg/logger"
	grpc2 "Webook/user/grpc"
	"github.com/spf13/viper"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
)

func InitGRPCxServer(userServer *grpc2.UserServiceServer,
	client *clientv3.Client, l logger.Logger) *grpcx.Server {
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
	userServer.Register(server)
	return &grpcx.Server{
		Server:     server,
		Port:       cfg.Port,
		Name:       "user",
		L:          l,
		EtcdClient: client,
		EtcdTTL:    cfg.EtcdTTL}
}
