package grpc2

import (
	searchv1 "Webook/api/proto/gen/api/proto/search/v1"
	"Webook/search/domain"
	"Webook/search/service"
	"context"
	"google.golang.org/grpc"
)

type SyncServiceServer struct {
	searchv1.UnimplementedSyncServiceServer
	syncSvc service.SyncService
}

func (s *SyncServiceServer) InputUser(ctx context.Context, request *searchv1.InputUserRequest) (*searchv1.InputUserResponse, error) {
	err := s.syncSvc.InputUser(ctx, s.toDomainUser(request.GetUser()))
	return &searchv1.InputUserResponse{}, err
}

func (s *SyncServiceServer) InputArticle(ctx context.Context, request *searchv1.InputArticleRequest) (*searchv1.InputArticleResponse, error) {
	err := s.syncSvc.InputArticle(ctx, s.toDomainArticle(request.GetArticle()))
	return &searchv1.InputArticleResponse{}, err
}

func (s *SyncServiceServer) InputAny(ctx context.Context, req *searchv1.InputAnyRequest) (*searchv1.InputAnyResponse, error) {
	err := s.syncSvc.InputAny(ctx, req.IndexName, req.DocId, req.Data)
	return &searchv1.InputAnyResponse{}, err
}

func NewSyncServiceServer(syncSvc service.SyncService) *SyncServiceServer {
	return &SyncServiceServer{
		syncSvc: syncSvc,
	}
}

func (s *SyncServiceServer) Register(server grpc.ServiceRegistrar) {
	searchv1.RegisterSyncServiceServer(server, s)
}

func (s *SyncServiceServer) toDomainUser(vuser *searchv1.User) domain.User {
	return domain.User{
		Id:       vuser.Id,
		Email:    vuser.Email,
		Nickname: vuser.Nickname,
	}
}

func (s *SyncServiceServer) toDomainArticle(art *searchv1.Article) domain.Article {
	return domain.Article{
		Id:      art.Id,
		Title:   art.Title,
		Status:  art.Status,
		Content: art.Content,
		Tags:    art.Tags,
	}
}
