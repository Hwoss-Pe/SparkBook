package ioc

import (
	pmtv1 "Webook/api/proto/gen/api/proto/payment/v1"
	"github.com/spf13/viper"
	etcdv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/resolver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func InitPaymentClient(etcdClient *etcdv3.Client) pmtv1.WechatPaymentServiceClient {
	type Config struct {
		Target string `json:"target"`
	}
	var cfg Config
	err := viper.UnmarshalKey("grpc.client.payment", &cfg)
	if err != nil {
		panic(err)
	}
	builder, err := resolver.NewBuilder(etcdClient)
	if err != nil {
		panic(err)
	}
	options := []grpc.DialOption{grpc.WithResolvers(builder), grpc.WithTransportCredentials(insecure.NewCredentials())}
	cc, err := grpc.Dial(cfg.Target, options...)
	if err != nil {
		panic(err)
	}
	return pmtv1.NewWechatPaymentServiceClient(cc)
}
