package client

import (
	intrv1 "Webook/api/proto/gen/api/proto/intr/v1"
	"Webook/interactive/service"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type InteractiveLocalAdapter struct {
	svc service.InteractiveService
}

func (i *InteractiveLocalAdapter) IncrReadCnt(ctx context.Context, in *intrv1.IncrReadCntRequest, opts ...grpc.CallOption) (*intrv1.IncrReadCntResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (i *InteractiveLocalAdapter) Like(ctx context.Context, in *intrv1.LikeRequest, opts ...grpc.CallOption) (*intrv1.LikeResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (i *InteractiveLocalAdapter) CancelLike(ctx context.Context, in *intrv1.CancelLikeRequest, opts ...grpc.CallOption) (*intrv1.CancelLikeResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (i *InteractiveLocalAdapter) Collect(ctx context.Context, in *intrv1.CollectRequest, opts ...grpc.CallOption) (*intrv1.CollectResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (i *InteractiveLocalAdapter) Get(ctx context.Context, in *intrv1.GetRequest, opts ...grpc.CallOption) (*intrv1.GetResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (i *InteractiveLocalAdapter) GetByIds(ctx context.Context, in *intrv1.GetByIdsRequest, opts ...grpc.CallOption) (*intrv1.GetByIdsResponse, error) {
	//TODO implement me
	panic("implement me")
}
