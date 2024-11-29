package grpc2

import (
	searchv1 "Webook/api/proto/gen/api/proto/search/v1"
	"Webook/search/domain"
	"Webook/search/service"
	"context"
	"github.com/ecodeclub/ekit/slice"
	"google.golang.org/grpc"
)

type SearchServiceServer struct {
	svc service.SearchService
	searchv1.UnimplementedSearchServiceServer
}

func NewSearchServiceServer(svc service.SearchService) *SearchServiceServer {
	return &SearchServiceServer{svc: svc}
}
func (s *SearchServiceServer) Register(server grpc.ServiceRegistrar) {
	searchv1.RegisterSearchServiceServer(server, s)
}
func (s *SearchServiceServer) Search(ctx context.Context, request *searchv1.SearchRequest) (*searchv1.SearchResponse, error) {
	resp, err := s.svc.Search(ctx, request.Uid, request.Expression)
	if err != nil {
		return nil, err
	}
	return &searchv1.SearchResponse{
		User: &searchv1.UserResult{
			Users: slice.Map(resp.Users, func(idx int, src domain.User) *searchv1.User {
				return &searchv1.User{
					Id:       src.Id,
					Nickname: src.Nickname,
					Email:    src.Email,
					Phone:    src.Phone,
				}
			}),
		},
		Article: &searchv1.ArticleResult{
			Articles: slice.Map(resp.Articles, func(idx int, src domain.Article) *searchv1.Article {
				return &searchv1.Article{
					Id:      src.Id,
					Title:   src.Title,
					Status:  src.Status,
					Content: src.Content,
				}
			}),
		},
	}, nil
}
