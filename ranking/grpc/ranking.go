package grpc2

import (
	rankingv1 "Webook/api/proto/gen/api/proto/ranking/v1"
	"Webook/ranking/domain"
	"Webook/ranking/service"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type RankingServiceServer struct {
	svc service.RankingService
	rankingv1.UnimplementedRankingServiceServer
}

func NewRankingServiceServer(svc service.RankingService) *RankingServiceServer {
	return &RankingServiceServer{
		svc: svc,
	}
}
func (r *RankingServiceServer) Register(server *grpc.Server) {
	rankingv1.RegisterRankingServiceServer(server, r)
}
func (r *RankingServiceServer) RankTopN(ctx context.Context, request *rankingv1.RankTopNRequest) (*rankingv1.RankTopNResponse, error) {
	err := r.svc.RankTopN(ctx)
	return &rankingv1.RankTopNResponse{}, err
}

func (r *RankingServiceServer) TopN(ctx context.Context, request *rankingv1.TopNRequest) (*rankingv1.TopNResponse, error) {
	domainArticles, err := r.svc.TopN(ctx)
	if err != nil {
		return &rankingv1.TopNResponse{}, err
	}
	res := make([]*rankingv1.Article, 0, len(domainArticles))
	for _, art := range domainArticles {
		res = append(res, convertToV(art))
	}
	return &rankingv1.TopNResponse{
		Articles: res,
	}, nil
}

func convertToV(da domain.Article) *rankingv1.Article {
	return &rankingv1.Article{
		Id:      da.Id,
		Title:   da.Title,
		Content: da.Content,
		Status:  int32(da.Status),
		Author: &rankingv1.Author{
			Id:   da.Author.Id,
			Name: da.Author.Name,
		},
		Ctime: timestamppb.New(da.Ctime),
		Utime: timestamppb.New(da.Utime),
	}
}
