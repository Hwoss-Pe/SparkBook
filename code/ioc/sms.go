package ioc

import (
	smsv1 "Webook/api/proto/gen/api/proto/sms/v1"
	"github.com/spf13/viper"
	etcdv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/resolver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func InitSmsRpcClient() smsv1.SmsServiceClient {
	type config struct {
		Target string `yaml:"target"`
	}
	var cfg config
	err := viper.UnmarshalKey("grpc.client.sms", &cfg)
	if err != nil {
		panic(err)
	}
	// 初始化 etcd 解析器，支持 "etcd:///service/sms" 形式的目标地址
	var opts []grpc.DialOption
	// 加载 etcd 客户端配置
	var etcdCfg etcdv3.Config
	if err := viper.UnmarshalKey("etcd", &etcdCfg); err == nil {
		if ecli, eErr := etcdv3.New(etcdCfg); eErr == nil {
			if builder, bErr := resolver.NewBuilder(ecli); bErr == nil {
				opts = append(opts, grpc.WithResolvers(builder))
			}
		}
	}
	// 关闭 TLS
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	conn, err := grpc.Dial(cfg.Target, opts...)
	if err != nil {
		panic(err)
	}
	return smsv1.NewSmsServiceClient(conn)
}
