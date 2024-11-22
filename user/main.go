package main

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func main() {
	initViperAndWatch()
	app := Init()
	err := app.GRPCServer.Serve()
	if err != nil {
		panic(err)
	}
}

func initViperAndWatch() {
	//	读取配置文件
	pfile := pflag.String("config", "config/dev.yaml",
		"配置文件路径")
	pflag.Parse()
	viper.SetConfigFile(*pfile)
	viper.WatchConfig()
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}
