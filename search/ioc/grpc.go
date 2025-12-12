package ioc

import (
	"Webook/pkg/grpcx"
	"Webook/pkg/logger"
	grpc2 "Webook/search/grpc"

	"github.com/spf13/viper"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
)

func InitGRPCxServer(syncRpc *grpc2.SyncServiceServer,
	searchRpc *grpc2.SearchServiceServer,
	ecli *clientv3.Client,
	l logger.Logger) *grpcx.Server {
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
	syncRpc.Register(server)
	searchRpc.Register(server)
	return &grpcx.Server{
		Server:     server,
		Port:       cfg.Port,
		Name:       "search",
		L:          l,
		EtcdTTL:    cfg.EtcdTTL,
		EtcdClient: ecli,
	}
}
