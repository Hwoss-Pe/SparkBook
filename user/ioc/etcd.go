package ioc

import (
	"github.com/spf13/viper"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func InitEtcdClient() *clientv3.Client {
	var cfg clientv3.Config
	//在main函数那里显式设置配置文件的路径
	err := viper.UnmarshalKey("etcd", &cfg)
	if err != nil {
		panic(err)
	}
	client, err := clientv3.New(cfg)
	if err != nil {
		panic(err)
	}
	return client
}
