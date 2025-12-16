package main

import (
	"Webook/pkg/grpcx"
	"Webook/pkg/saramax"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
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
	for _, c := range app.consumers {
		if err := c.Start(); err != nil {
			panic(err)
		}
	}
	err := app.server.Serve()
	if err != nil {
		panic(err)
	}
}

type App struct {
	server    *grpcx.Server
	consumers []saramax.Consumer
}
