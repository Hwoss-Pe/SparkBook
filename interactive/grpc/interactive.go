package grpc2

import (
	"Webook/api/proto/gen/api/proto/intr/v1"
	"Webook/interactive/domain"
	"Webook/interactive/service"
	"context"

	"google.golang.org/grpc"
)

type InteractiveServiceServer struct {
	intrv1.UnimplementedInteractiveServiceServer
	svc   service.InteractiveService
	notif service.NotificationService
}

func (i *InteractiveServiceServer) Register(server grpc.ServiceRegistrar) {
	intrv1.RegisterInteractiveServiceServer(server, i)
}

func NewInteractiveServiceServer(svc service.InteractiveService, notif service.NotificationService) *InteractiveServiceServer {
	return &InteractiveServiceServer{svc: svc, notif: notif}
}

func (i *InteractiveServiceServer) IncrReadCnt(ctx context.Context, request *intrv1.IncrReadCntRequest) (*intrv1.IncrReadCntResponse, error) {
	err := i.svc.IncrReadCnt(ctx, request.Biz, request.BizId)
	return &intrv1.IncrReadCntResponse{}, err
}

func (i *InteractiveServiceServer) Like(ctx context.Context, request *intrv1.LikeRequest) (*intrv1.LikeResponse, error) {
	err := i.svc.Like(ctx, request.Biz, request.BizId, request.Uid)
	return &intrv1.LikeResponse{}, err
}

func (i *InteractiveServiceServer) CancelLike(ctx context.Context, request *intrv1.CancelLikeRequest) (*intrv1.CancelLikeResponse, error) {
	err := i.svc.CancelLike(ctx, request.Biz, request.BizId, request.Uid)
	return &intrv1.CancelLikeResponse{}, err
}

func (i *InteractiveServiceServer) Collect(ctx context.Context, request *intrv1.CollectRequest) (*intrv1.CollectResponse, error) {
	err := i.svc.Collect(ctx, request.Biz, request.BizId, request.Cid, request.Uid)
	return &intrv1.CollectResponse{}, err
}

func (i *InteractiveServiceServer) CancelCollect(ctx context.Context, request *intrv1.CancelCollectRequest) (*intrv1.CancelCollectResponse, error) {
	err := i.svc.CancelCollect(ctx, request.Biz, request.BizId, request.Cid, request.Uid)
	return &intrv1.CancelCollectResponse{}, err
}

func (i *InteractiveServiceServer) Get(ctx context.Context, request *intrv1.GetRequest) (*intrv1.GetResponse, error) {
	res, err := i.svc.Get(ctx, request.Biz, request.BizId, request.Uid)
	if err != nil {
		return nil, err
	}
	return &intrv1.GetResponse{
		Intr: i.toDTO(res),
	}, nil
}

func (i *InteractiveServiceServer) GetByIds(ctx context.Context, request *intrv1.GetByIdsRequest) (*intrv1.GetByIdsResponse, error) {
	if len(request.Ids) == 0 {
		return &intrv1.GetByIdsResponse{}, nil
	}
	data, err := i.svc.GetByIds(ctx, request.Biz, request.Ids)
	if err != nil {
		return nil, err
	}

	res := make(map[int64]*intrv1.Interactive, len(data))
	for k, v := range data {
		res[k] = i.toDTO(v)
	}
	return &intrv1.GetByIdsResponse{
		Intrs: res,
	}, nil
}

func (i *InteractiveServiceServer) GetCollectedBizIds(ctx context.Context, request *intrv1.GetCollectedBizIdsRequest) (*intrv1.GetCollectedBizIdsResponse, error) {
	bizIds, total, err := i.svc.GetCollectedBizIds(ctx, request.Biz, request.Uid, int(request.Offset), int(request.Limit))
	if err != nil {
		return nil, err
	}
	return &intrv1.GetCollectedBizIdsResponse{
		BizIds: bizIds,
		Total:  total,
	}, nil
}

func (i *InteractiveServiceServer) GetNotifications(ctx context.Context, req *intrv1.GetNotificationsRequest) (*intrv1.GetNotificationsResponse, error) {
	return i.notif.GetNotifications(ctx, req.GetUid(), req.GetType(), int(req.GetOffset()), int(req.GetLimit()))
}

func (i *InteractiveServiceServer) MarkRead(ctx context.Context, req *intrv1.MarkReadRequest) (*intrv1.MarkReadResponse, error) {
	err := i.notif.MarkRead(ctx, req.GetUid(), req.GetIds(), req.GetType())
	return &intrv1.MarkReadResponse{}, err
}

func (i *InteractiveServiceServer) GetUnreadCounts(ctx context.Context, req *intrv1.GetUnreadCountsRequest) (*intrv1.GetUnreadCountsResponse, error) {
	return i.notif.GetUnreadCounts(ctx, req.GetUid())
}

func (i *InteractiveServiceServer) toDTO(intr domain.Interactive) *intrv1.Interactive {
	return &intrv1.Interactive{
		Biz:        intr.Biz,
		BizId:      intr.BizId,
		ReadCnt:    intr.ReadCnt,
		LikeCnt:    intr.LikeCnt,
		CollectCnt: intr.CollectCnt,
		Liked:      intr.Liked,
		Collected:  intr.Collected,
	}
}
