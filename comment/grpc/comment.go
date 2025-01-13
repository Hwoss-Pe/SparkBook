package grpc2

import (
	commentv1 "Webook/api/proto/gen/api/proto/comment/v1"
	"Webook/comment/domain"
	"Webook/comment/service"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
	"math"
)

type CommentServiceServer struct {
	svc service.CommentService

	commentv1.UnimplementedCommentServiceServer
}

func NewGrpcServer(svc service.CommentService) *CommentServiceServer {
	return &CommentServiceServer{
		svc: svc,
	}
}
func (c *CommentServiceServer) Register(server grpc.ServiceRegistrar) {
	commentv1.RegisterCommentServiceServer(server, c)
}
func (c *CommentServiceServer) GetCommentList(ctx context.Context, request *commentv1.CommentListRequest) (*commentv1.CommentListResponse, error) {
	minID := request.MinId
	if minID <= 0 {
		minID = math.MaxInt64
	}
	list, err := c.svc.GetCommentList(ctx, request.Biz, request.Bizid, minID, request.Limit)
	if err != nil {
		return nil, err
	}
	return &commentv1.CommentListResponse{
		Comments: c.toDTO(list),
	}, nil
}

func (c *CommentServiceServer) DeleteComment(ctx context.Context, request *commentv1.DeleteCommentRequest) (*commentv1.DeleteCommentResponse, error) {
	err := c.svc.DeleteComment(ctx, request.Id)
	return &commentv1.DeleteCommentResponse{}, err
}

func (c *CommentServiceServer) CreateComment(ctx context.Context, request *commentv1.CreateCommentRequest) (*commentv1.CreateCommentResponse, error) {
	err := c.svc.CreateComment(ctx, c.toDomain(request.GetComment()))
	return &commentv1.CreateCommentResponse{}, err
}

func (c *CommentServiceServer) GetMoreReplies(ctx context.Context, request *commentv1.GetMoreRepliesRequest) (*commentv1.GetMoreRepliesResponse, error) {
	replies, err := c.svc.GetMoreReplies(ctx, request.Rid, request.MaxId, request.Limit)
	if err != nil {
		return nil, err
	}
	return &commentv1.GetMoreRepliesResponse{
		Replies: c.toDTO(replies),
	}, nil
}

// 这个方法主要是对返回一个数组，但是如果想要往上查找父评论的话需要走递归
func (c *CommentServiceServer) toDTO(domainComments []domain.Comment) []*commentv1.Comment {
	rpcComments := make([]*commentv1.Comment, 0, len(domainComments))
	//拿这个进行快速选择id进行填充
	rpcCommentMap := make(map[int64]*commentv1.Comment, len(domainComments))
	for _, domainComment := range domainComments {
		//需要一个目的映射的map
		rpcCommentMap[domainComment.Id] = &commentv1.Comment{
			Id:      domainComment.Id,
			Uid:     domainComment.Commentator.ID,
			Biz:     domainComment.Biz,
			Bizid:   domainComment.BizId,
			Content: domainComment.Content,
			Ctime:   timestamppb.New(domainComment.CTime),
			Utime:   timestamppb.New(domainComment.UTime),
		}
	}
	for _, domainComment := range domainComments {
		rpcComment := &commentv1.Comment{
			Id:      domainComment.Id,
			Uid:     domainComment.Commentator.ID,
			Biz:     domainComment.Biz,
			Bizid:   domainComment.BizId,
			Content: domainComment.Content,
			Ctime:   timestamppb.New(domainComment.CTime),
			Utime:   timestamppb.New(domainComment.UTime),
		}
		if domainComment.RootComment != nil {
			rpcComment.RootComment = rpcCommentMap[domainComment.RootComment.Id]
		}
		if domainComment.ParentComment != nil {
			rpcComment.ParentComment = rpcCommentMap[domainComment.ParentComment.Id]
		}
		rpcComments = append(rpcComments, rpcComment)
	}
	return rpcComments

}

func (c *CommentServiceServer) toDomain(comment *commentv1.Comment) domain.Comment {
	domainComment := domain.Comment{
		Id:      comment.Id,
		Biz:     comment.Biz,
		BizId:   comment.Bizid,
		Content: comment.Content,
		Commentator: domain.User{
			ID: comment.Uid,
		},
	}
	if comment.ParentComment != nil {
		domainComment.ParentComment = &domain.Comment{
			Id: comment.GetParentComment().GetId(),
		}
	}
	if comment.GetRootComment() != nil {
		domainComment.RootComment = &domain.Comment{
			Id: comment.GetRootComment().GetId(),
		}
	}
	return domainComment
}
