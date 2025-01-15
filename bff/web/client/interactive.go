package client

import (
	intrv1 "Webook/api/proto/gen/api/proto/intr/v1"
	"github.com/ecodeclub/ekit/syncx/atomicx"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"math/rand"
)

type InteractiveClient struct {
	remote intrv1.InteractiveServiceClient
	local  *InteractiveLocalAdapter

	//	流量控制[0,100)
	threshold *atomicx.Value[int32]
}

func (i *InteractiveClient) IncrReadCnt(ctx context.Context, in *intrv1.IncrReadCntRequest, opts ...grpc.CallOption) (*intrv1.IncrReadCntResponse, error) {
	return i.selectClient().IncrReadCnt(ctx, in)
}

func (i *InteractiveClient) Like(ctx context.Context, in *intrv1.LikeRequest, opts ...grpc.CallOption) (*intrv1.LikeResponse, error) {
	return i.selectClient().Like(ctx, in)
}

func (i *InteractiveClient) CancelLike(ctx context.Context, in *intrv1.CancelLikeRequest, opts ...grpc.CallOption) (*intrv1.CancelLikeResponse, error) {
	return i.selectClient().CancelLike(ctx, in)
}

func (i *InteractiveClient) Collect(ctx context.Context, in *intrv1.CollectRequest, opts ...grpc.CallOption) (*intrv1.CollectResponse, error) {
	return i.selectClient().Collect(ctx, in)
}

func (i *InteractiveClient) Get(ctx context.Context, in *intrv1.GetRequest, opts ...grpc.CallOption) (*intrv1.GetResponse, error) {
	return i.selectClient().Get(ctx, in)
}

func (i *InteractiveClient) GetByIds(ctx context.Context, in *intrv1.GetByIdsRequest, opts ...grpc.CallOption) (*intrv1.GetByIdsResponse, error) {
	return i.selectClient().GetByIds(ctx, in)
}

func NewInteractiveClient(remote intrv1.InteractiveServiceClient,
	local *InteractiveLocalAdapter, threshold int32) *InteractiveClient {
	return &InteractiveClient{
		remote:    remote,
		local:     local,
		threshold: atomicx.NewValueOf(threshold),
	}

}

// 做一个限流的分发
func (i *InteractiveClient) selectClient() intrv1.InteractiveServiceClient {
	num := rand.Int31n(100)
	if num < i.threshold.Load() {
		return i.remote
	}
	return i.local
}

func (i *InteractiveClient) UpdateThreshold(threshold int32) {
	i.threshold.Store(threshold)
}
