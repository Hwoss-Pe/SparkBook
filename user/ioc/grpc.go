package ioc

import (
	"Webook/pkg/logger"
	"Webook/user/grpc"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func InitGRPCxServer(userServer *grpc.UserServiceServer,
	client clientv3.Client, l logger.Logger) {

}
