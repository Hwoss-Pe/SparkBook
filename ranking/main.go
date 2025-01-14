package main

import (
	"Webook/pkg/grpcx"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type App struct {
	server *grpcx.Server
}

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
	err := app.server.Serve()
	if err != nil {
		panic(err)
	}
}
