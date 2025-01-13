package main

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"net"
)

func initViperV2Watch() {
	cfile := pflag.String("config",
		"config/dev.yaml", "配置文件路径")
	pflag.Parse()
	// 直接指定文件路径
	viper.SetConfigFile(*cfile)
	viper.WatchConfig()
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}

func main() {
	initViperV2Watch()
	app := Init()
	//这里也可以去加服务注册发现
	listen, err := net.Listen("tcp", ":"+"8080")
	err = app.server.Server.Serve(listen)
	if err != nil {
		panic(err)
	}
}
