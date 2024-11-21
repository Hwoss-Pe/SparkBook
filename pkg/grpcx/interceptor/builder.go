package interceptor

import (
	"context"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"net"
	"strings"
)

type Builder struct {
}

// PeerName 获取对端应用名称
func (b *Builder) PeerName(ctx context.Context) string {
	return b.grpcHeaderValue(ctx, "app")
}

// PeerIP 获取对端ip
func (b *Builder) PeerIP(ctx context.Context) string {
	//ctx看看有没有设置
	clientIP := b.grpcHeaderValue(ctx, "client-ip")
	if clientIP != "" {
		return clientIP
	}
	//从grpc里取对端ip,正常来说是会有在x-forwarded-for里面，这里过网关后会在后面添加
	pr, ok := peer.FromContext(ctx)
	if !ok {
		return ""
	}
	if pr.Addr == net.Addr(nil) {
		return ""
	}
	addSlice := strings.Split(pr.Addr.String(), ":")
	//addSlice := strings.Split(pr.Addr.String(), ",")
	//只返回最开始的ip
	if len(addSlice) > 1 {
		return addSlice[0]
	}
	return ""
}

// grpcHeaderValue 获取grpc携带的请求头
func (b *Builder) grpcHeaderValue(ctx context.Context, key string) string {
	if key == "" {
		return ""
	}
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}
	return strings.Join(md.Get(key), ";")
}
