//go:build wireinject

package main

import (
	"Webook/pkg/wego"
	grpc2 "Webook/sms/grpc"
	"Webook/sms/ioc"
	"Webook/sms/repository"
	"Webook/sms/repository/dao"
	"Webook/sms/service"
	"Webook/sms/service/async"
	"Webook/sms/service/tencent"
	"github.com/google/wire"
)

// 提供 SmsTencentService 实现的 ProviderSet
var SmsTencentProviderSet = wire.NewSet(
	ioc.InitSmsTencentService,                                     // 提供 InitSmsTencentService
	wire.Bind(new(service.Service), new(*tencent.TencentService)), // 显式绑定接口和实现
)

// 提供异步短信服务的 ProviderSet
var AsyncServiceProviderSet = wire.NewSet(
	async.NewAsyncService, // 提供 NewAsyncService
	wire.Bind(new(service.Service), new(*async.AsyncService)), // 显式绑定接口和实现
)

// 提供数据库和 dao 的 ProviderSet
var DaoProviderSet = wire.NewSet(
	ioc.InitDB,             // 初始化 DB
	dao.NewGORMAsyncSmsDAO, // 创建 DAO 实现
)

// 提供服务存储库的 ProviderSet
var RepositoryProviderSet = wire.NewSet(
	repository.NewAsyncSMSRepository, // 创建异步短信服务的存储库
)

// 创建日志和配置的 ProviderSet
var LoggerProviderSet = wire.NewSet(
	ioc.InitLogger, // 初始化日志
)

// 提供 Etcd 客户端的 ProviderSet
var EtcdClientProviderSet = wire.NewSet(
	ioc.InitEtcdClient, // 初始化 Etcd 客户端
)

// 创建 GRPC 服务和服务端的 ProviderSet
var GRPCProviderSet = wire.NewSet(
	grpc2.NewSmsServiceServer, // 提供 GRPC 服务
	ioc.InitGRPCxServer,       // 初始化 GRPC 服务器
)

// InitializeApp 使用 wire 注入依赖并返回最终的 app 实例
func Init() *wego.App {
	wire.Build(
		DaoProviderSet,                           // 提供数据库和 DAO
		RepositoryProviderSet,                    // 提供服务存储库
		LoggerProviderSet,                        // 提供日志
		EtcdClientProviderSet,                    // 提供 Etcd 客户端
		SmsTencentProviderSet,                    // 提供短信服务
		AsyncServiceProviderSet,                  // 提供异步短信服务
		GRPCProviderSet,                          // 提供 GRPC 服务
		wire.Struct(new(wego.App), "GRPCServer"), // 最终构建 app
	)
	return nil
}
