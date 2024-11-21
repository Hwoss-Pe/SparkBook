package grpcx

import (
	"Webook/pkg/logger"
	"Webook/pkg/netx"
	"context"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/endpoints"
	"google.golang.org/grpc"
	"net"
	"strconv"
	"time"
)

type Server struct {
	*grpc.Server
	Port int
	//	ETCD 的续约机制
	EtcdTTL     int64
	EtcdClient  *clientv3.Client
	etcdManager endpoints.Manager
	etcdKey     string
	cancel      func()
	Name        string
	L           logger.Logger
}

func (s *Server) Serve() error {
	ctx, cancelFunc := context.WithCancel(context.Background())
	s.cancel = cancelFunc
	port := strconv.Itoa(s.Port)
	listen, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}
	//	启动成功后在进行注册
	err = s.register(ctx, port)
	if err != nil {
		return err
	}
	return s.Server.Serve(listen)
}
func (s *Server) register(ctx context.Context, port string) error {
	cli := s.EtcdClient
	serviceName := "service/" + s.Name
	//获取管理服务的注册和发现的对象
	manager, err := endpoints.NewManager(cli, serviceName)
	if err != nil {
		return err
	}
	s.etcdManager = manager
	//获取本机当前外网的ip
	ip := netx.GetOutboundIP()
	s.etcdKey = ip + ":" + port
	addr := ip + ":" + port
	//开始租约
	leaseGrantResponse, err := cli.Grant(ctx, s.EtcdTTL)
	//开启自动续约
	ch, err := cli.KeepAlive(ctx, leaseGrantResponse.ID)
	if err != nil {
		return err
	}
	go func() {
		//在这异步可以看到续约信息
		for chResp := range ch {
			s.L.Debug("续约：", logger.String("resp", chResp.String()))
		}
	}()

	//一定要在这里管理节点，并且要携带租约，这样续约过期就会自动剔除
	return manager.AddEndpoint(ctx, s.etcdKey,
		endpoints.Endpoint{
			Addr: addr,
		}, clientv3.WithLease(leaseGrantResponse.ID))
}
func (s *Server) close() error {
	s.cancel()
	if s.etcdManager != nil {
		ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second)
		defer cancelFunc()
		//单独开一个ctx来进行删除节点的操作
		err := s.etcdManager.DeleteEndpoint(ctx, s.etcdKey)
		if err != nil {
			return err
		}
	}
	err := s.EtcdClient.Close()
	if err != nil {
		return err
	}
	//平滑关闭一般0.5s
	s.Server.GracefulStop()
	return nil
}
