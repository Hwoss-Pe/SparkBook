package ioc

import (
	grpc2 "Webook/cronjob/grpc"
	"Webook/pkg/grpcx"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func InitGRPCxServer(cronJobGrpc *grpc2.CronJobServiceServer) *grpcx.Server {
	type Config struct {
		Addr int `yaml:"addr"`
	}
	var cfg Config
	err := viper.UnmarshalKey("grpc.server", &cfg)
	if err != nil {
		panic(err)
	}
	server := grpc.NewServer()
	cronJobGrpc.Register(server)
	return &grpcx.Server{
		Server: server,
		Port:   cfg.Addr,
	}
}
