package main

import (
	"Webook/pkg/ginx"
	"Webook/pkg/grpcx"
	"Webook/pkg/saramax"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func main() {
	initViper()
	initPrometheus()
	app := Init()
	for _, consumer := range app.consumers {
		err := consumer.Start()
		if err != nil {
			panic(err)
		}
	}
	go func() {
		err := app.migratorServer.Start()
		if err != nil {
			panic(err)
		}
	}()
	err := app.server.Serve()
	if err != nil {
		panic(err)
	}
}

func initViper() {
	cfile := pflag.String("config", "config/dev.yaml", "文件路径")
	pflag.Parse()
	viper.SetConfigFile(*cfile)
	viper.WatchConfig()
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}

type App struct {
	server         *grpcx.Server
	migratorServer *ginx.Server
	consumers      []saramax.Consumer
}
