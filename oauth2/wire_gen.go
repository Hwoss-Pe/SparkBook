// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"Webook/oauth2/grpc"
	"Webook/oauth2/ioc"
	"github.com/google/wire"
)

// Injectors from wire.go:

func Init() *App {
	logger := ioc.InitLogger()
	service := ioc.InitPrometheus(logger)
	oauth2ServiceServer := grpc.NewOauth2ServiceServer(service)
	client := ioc.InitEtcdClient()
	server := ioc.InitGRPCxServer(oauth2ServiceServer, client, logger)
	app := &App{
		server: server,
	}
	return app
}

// wire.go:

var thirdProvider = wire.NewSet(ioc.InitLogger)
