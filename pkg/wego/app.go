package wego

import (
	"Webook/pkg/ginx"
	"Webook/pkg/grpcx"
	"Webook/pkg/saramax"
)

// App 是作为wire聚合的时候使用这个结构体，也不是所有服务都需要下面的字段
type App struct {
	GRPCServer *grpcx.Server
	WebServer  *ginx.Server
	Consumers  []saramax.Consumer
}
