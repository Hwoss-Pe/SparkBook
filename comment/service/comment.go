package service

import (
	"Webook/comment/domain"
	"Webook/comment/events"
	"Webook/comment/repository"
	"golang.org/x/net/context"
)

type CommentService interface {
	// GetCommentList Comment的id为0 获取一级评论
	// 按照 ID 倒序排序
	GetCommentList(ctx context.Context, biz string, bizId, minID, limit int64) ([]domain.Comment, error)
	// DeleteComment 删除评论，删除本评论何其子评论
	DeleteComment(ctx context.Context, id int64) error
	// CreateComment 创建评论
	CreateComment(ctx context.Context, comment domain.Comment) error
	GetMoreReplies(ctx context.Context, rid int64, maxID int64, limit int64) ([]domain.Comment, error)
}
type commentService struct {
	repo     repository.CommentRepository
	producer events.Producer
}

func (c *commentService) GetCommentList(ctx context.Context, biz string, bizId, minID, limit int64) ([]domain.Comment, error) {
	comments, err := c.repo.FindByBiz(ctx, biz, bizId, minID, limit)
	if err != nil {
		return nil, err
	}
	return comments, nil
}

func (c *commentService) DeleteComment(ctx context.Context, id int64) error {
	if c.producer != nil {
		return c.producer.ProduceDelete(events.CommentDeleteEvent{Id: id})
	}
	return c.repo.DeleteComment(ctx, domain.Comment{Id: id})
}

func (c *commentService) CreateComment(ctx context.Context, comment domain.Comment) error {
	if c.producer != nil {
		var pid int64
		var rid int64
		if comment.ParentComment != nil {
			pid = comment.ParentComment.Id
		}
		if comment.RootComment != nil {
			rid = comment.RootComment.Id
		}
		return c.producer.ProduceCreate(events.CommentCreateEvent{
			Biz: comment.Biz, BizId: comment.BizId, Uid: comment.Commentator.ID, Content: comment.Content, ParentId: pid, RootId: rid,
		})
	}
	return c.repo.CreateComment(ctx, comment)
}

func (c *commentService) GetMoreReplies(ctx context.Context, rid int64, maxID int64, limit int64) ([]domain.Comment, error) {
	return c.repo.GetMoreReplies(ctx, rid, maxID, limit)
}

func NewCommentSvc(repo repository.CommentRepository, p events.Producer) CommentService {
	return &commentService{repo: repo, producer: p}
}
