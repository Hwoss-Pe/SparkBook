package main

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func main() {
	initViperV2Watch()
	app := Init()
	err := app.GRPCServer.Serve()
	if err != nil {
		panic(err)
	}
}
func initViperV2Watch() {
	cfile := pflag.String("config", "config/dev.yaml", "xx")
	pflag.Parse()
	viper.SetConfigFile(*cfile)
	viper.WatchConfig()
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}
