package repository

import (
	"Webook/comment/domain"
	"Webook/comment/repository/dao"
	"Webook/pkg/logger"
	"database/sql"
	"time"

	"golang.org/x/net/context"
	"golang.org/x/sync/errgroup"
)

type CommentRepository interface {
	//FindByBiz 根据id倒序查找，并且返回每个评论的三条直接回复
	FindByBiz(ctx context.Context, biz string, bizId, minID, limit int64) ([]domain.Comment, error)
	DeleteComment(ctx context.Context, comment domain.Comment) error
	CreateComment(ctx context.Context, comment domain.Comment) error
	//	GetCommentByIds，支持单条获取
	GetCommentByIds(ctx context.Context, id []int64) ([]domain.Comment, error)
	GetMoreReplies(ctx context.Context, rid int64, id int64, limit int64) ([]domain.Comment, error)
}

type CachedCommentRepo struct {
	dao dao.CommentDAO
	l   logger.Logger
}

func NewCommentRepo(commentDAO dao.CommentDAO, l logger.Logger) CommentRepository {
	return &CachedCommentRepo{
		dao: commentDAO,
		l:   l,
	}
}

func (c *CachedCommentRepo) FindByBiz(ctx context.Context, biz string, bizId, minID, limit int64) ([]domain.Comment, error) {
	daoComments, err := c.dao.FindByBiz(ctx, biz, bizId, minID, limit)
	if err != nil {
		return nil, err
	}
	res := make([]domain.Comment, len(daoComments))
	downgraded := ctx.Value("downgraded") == "true"
	var eg errgroup.Group
	for i, cmt := range daoComments {
		comment := cmt
		cm := c.toDomain(comment)
		res[i] = cm
		if downgraded {
			continue
		}
		idx := i
		eg.Go(func() error {
			rs, err := c.dao.FindRepliesByPid(ctx, comment.Id, 0, 3)
			if err != nil {
				c.l.Error("查询子评论失败", logger.Error(err))
				return nil
			}
			children := make([]domain.Comment, 0, len(rs))
			for _, r := range rs {
				children = append(children, c.toDomain(r))
			}
			res[idx].Children = children
			return nil
		})
	}
	return res, eg.Wait()
}

func (c *CachedCommentRepo) DeleteComment(ctx context.Context, comment domain.Comment) error {
	return c.dao.Delete(ctx, dao.Comment{
		Id: comment.Id,
	})
}

func (c *CachedCommentRepo) CreateComment(ctx context.Context, comment domain.Comment) error {
	return c.dao.Insert(ctx, c.toEntity(comment))
}

func (c *CachedCommentRepo) GetCommentByIds(ctx context.Context, ids []int64) ([]domain.Comment, error) {
	vals, err := c.dao.FindOneByIDs(ctx, ids)
	if err != nil {
		return nil, err
	}
	comments := make([]domain.Comment, 0, len(vals))

	for _, val := range vals {
		comment := c.toDomain(val)
		comments = append(comments, comment)
	}
	return comments, nil
}

func (c *CachedCommentRepo) GetMoreReplies(ctx context.Context, rid int64, maxID int64, limit int64) ([]domain.Comment, error) {
	cs, err := c.dao.FindRepliesByRid(ctx, rid, maxID, limit)
	if err != nil {
		return nil, err
	}
	res := make([]domain.Comment, 0, len(cs))
	for _, cm := range cs {
		res = append(res, c.toDomain(cm))
	}
	return res, nil
}

func (c *CachedCommentRepo) toDomain(dc dao.Comment) domain.Comment {
	val := domain.Comment{
		Id: dc.Id,
		Commentator: domain.User{
			ID: dc.Uid,
		},
		Biz:     dc.Biz,
		BizId:   dc.BizID,
		Content: dc.Content,
		CTime:   time.UnixMilli(dc.Ctime),
		UTime:   time.UnixMilli(dc.Utime),
	}
	if dc.RootID.Valid {
		val.RootComment = &domain.Comment{
			Id: dc.RootID.Int64,
		}
	}
	if dc.PID.Valid {
		val.ParentComment = &domain.Comment{
			Id: dc.PID.Int64,
		}
		if dc.ParentComment != nil {
			val.ParentComment.Commentator = domain.User{ID: dc.ParentComment.Uid}
		}
	}

	return val
}

func (c *CachedCommentRepo) toEntity(domainComment domain.Comment) dao.Comment {
	daoComment := dao.Comment{
		Id:      domainComment.Id,
		Uid:     domainComment.Commentator.ID,
		Biz:     domainComment.Biz,
		BizID:   domainComment.BizId,
		Content: domainComment.Content,
	}
	if domainComment.RootComment != nil {
		daoComment.RootID = sql.NullInt64{
			Valid: true,
			Int64: domainComment.RootComment.Id,
		}
	}
	if domainComment.ParentComment != nil {
		daoComment.PID = sql.NullInt64{
			Valid: true,
			Int64: domainComment.ParentComment.Id,
		}
	}
	daoComment.Ctime = time.Now().UnixMilli()
	daoComment.Utime = time.Now().UnixMilli()
	return daoComment
}
