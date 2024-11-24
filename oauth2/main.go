package main

import (
	"Webook/pkg/grpcx"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type App struct {
	server *grpcx.Server
}

func initViperWatch() {
	cflag := pflag.String("config", "config/dev.yaml", "config")
	pflag.Parse()
	viper.SetConfigFile(*cflag)
	viper.WatchConfig()
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}

func main() {
	initViperWatch()
	app := Init()
	err := app.server.Serve()
	if err != nil {
		panic(err)
	}
}
