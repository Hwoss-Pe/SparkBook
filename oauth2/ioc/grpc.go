package ioc

import (
	grpc2 "Webook/oauth2/grpc"
	"Webook/pkg/grpcx"
	"Webook/pkg/logger"
	"github.com/spf13/viper"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
)

func InitGRPCxServer(oauth2Server *grpc2.Oauth2ServiceServer, cli *clientv3.Client, l logger.Logger) *grpcx.Server {
	type Config struct {
		Addr    string `yaml:"addr"`
		EtcdTTL int64  `yaml:"etcdTTL"`
	}
	var cfg Config
	err := viper.UnmarshalKey("grpc.server", &cfg)
	if err != nil {
		panic(err)
	}
	server := grpc.NewServer()
	oauth2Server.Register(server)
	return &grpcx.Server{
		Server:     server,
		EtcdClient: cli,
		Name:       "auth",
		EtcdTTL:    cfg.EtcdTTL,
		L:          l,
	}
}
