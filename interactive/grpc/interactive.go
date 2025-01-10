package grpc2

import (
	"Webook/api/proto/gen/api/proto/intr/v1"
	"Webook/interactive/service"
	"context"
	"google.golang.org/grpc"
)

type InteractiveServiceServer struct {
	intrv1.UnimplementedInteractiveServiceServer
	svc service.InteractiveService
}

func (i *InteractiveServiceServer) Register(server grpc.ServiceRegistrar) {
	intrv1.RegisterInteractiveServiceServer(server, i)
}

func NewInteractiveServiceServer(svc service.InteractiveService) *InteractiveServiceServer {
	return &InteractiveServiceServer{svc: svc}
}

func (i *InteractiveServiceServer) IncrReadCnt(ctx context.Context, request *intrv1.IncrReadCntRequest) (*intrv1.IncrReadCntResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (i *InteractiveServiceServer) Like(ctx context.Context, request *intrv1.LikeRequest) (*intrv1.LikeResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (i *InteractiveServiceServer) CancelLike(ctx context.Context, request *intrv1.CancelLikeRequest) (*intrv1.CancelLikeResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (i *InteractiveServiceServer) Collect(ctx context.Context, request *intrv1.CollectRequest) (*intrv1.CollectResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (i *InteractiveServiceServer) Get(ctx context.Context, request *intrv1.GetRequest) (*intrv1.GetResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (i *InteractiveServiceServer) GetByIds(ctx context.Context, request *intrv1.GetByIdsRequest) (*intrv1.GetByIdsResponse, error) {
	//TODO implement me
	panic("implement me")
}
