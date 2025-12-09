package web

import (
	commentv1 "Webook/api/proto/gen/api/proto/comment/v1"
	userv1 "Webook/api/proto/gen/api/proto/user/v1"
	"Webook/pkg/ginx"
	"Webook/pkg/logger"
	"time"

	"github.com/gin-gonic/gin"
)

type CommentHandler struct {
	client  commentv1.CommentServiceClient
	userSvc userv1.UsersServiceClient
	l       logger.Logger
}

func NewCommentHandler(client commentv1.CommentServiceClient, userSvc userv1.UsersServiceClient, l logger.Logger) *CommentHandler {
	return &CommentHandler{client: client, userSvc: userSvc, l: l}
}

func (h *CommentHandler) RegisterRoute(s *gin.Engine) {
	g := s.Group("/comment")
	g.GET("/list", ginx.WrapReq[CommentListReq](h.List))
	g.POST("/create", ginx.WrapClaimsAndReq[CreateReq](h.Create))
	g.POST("/delete", ginx.WrapClaimsAndReq[DeleteReq](h.Delete))
	g.GET("/replies", ginx.WrapReq[RepliesReq](h.Replies))
}

type CommentListReq struct {
	Biz   string `form:"biz"`
	BizId int64  `form:"bizid"`
	MinId int64  `form:"min_id"`
	Limit int64  `form:"limit"`
}

type DeleteReq struct {
	Id int64 `json:"id"`
}

type InnerId struct {
	Id int64 `json:"id"`
}

type CreateReq struct {
	Comment struct {
		Uid           int64    `json:"uid"`
		Biz           string   `json:"biz"`
		BizId         int64    `json:"bizid"`
		Content       string   `json:"content"`
		RootComment   *InnerId `json:"root_comment"`
		ParentComment *InnerId `json:"parent_comment"`
	} `json:"comment"`
}

type RepliesReq struct {
	Rid   int64 `form:"rid"`
	MaxId int64 `form:"max_id"`
	Limit int64 `form:"limit"`
}

func (h *CommentHandler) List(ctx *gin.Context, req CommentListReq) (ginx.Result, error) {
	resp, err := h.client.GetCommentList(ctx, &commentv1.CommentListRequest{
		Biz:   req.Biz,
		Bizid: req.BizId,
		MinId: req.MinId,
		Limit: req.Limit,
	})
	if err != nil {
		return ginx.Result{Code: 5, Msg: "系统错误"}, err
	}
	cs := resp.GetComments()
	uids := make(map[int64]struct{}, len(cs))
	for _, c := range cs {
		if c.GetUid() > 0 {
			uids[c.GetUid()] = struct{}{}
		}
		if pc := c.GetParentComment(); pc != nil && pc.GetUid() > 0 {
			uids[pc.GetUid()] = struct{}{}
		}
	}
	type ui struct{ name, avatar string }
	profiles := make(map[int64]ui, len(uids))
	for uid := range uids {
		pr, er := h.userSvc.Profile(ctx, &userv1.ProfileRequest{Id: uid})
		if er != nil || pr == nil || pr.User == nil {
			continue
		}
		profiles[uid] = ui{name: pr.User.Nickname, avatar: pr.User.Avatar}
	}
	roots := make(map[int64]map[string]any, len(cs))
	order := make([]int64, 0, len(cs))
	for _, c := range cs {
		if c.GetRootComment() == nil {
			item := map[string]any{
				"id":       c.GetId(),
				"uid":      c.GetUid(),
				"name":     profiles[c.GetUid()].name,
				"avatar":   profiles[c.GetUid()].avatar,
				"biz":      c.GetBiz(),
				"bizid":    c.GetBizid(),
				"content":  c.GetContent(),
				"ctime":    c.GetCtime().AsTime().Format(time.DateTime),
				"utime":    c.GetUtime().AsTime().Format(time.DateTime),
				"children": make([]map[string]any, 0, 3),
			}
			roots[c.GetId()] = item
			order = append(order, c.GetId())
		}
	}
	for _, c := range cs {
		if rc := c.GetRootComment(); rc != nil {
			rid := rc.GetId()
			root := roots[rid]
			if root == nil {
				continue
			}
			child := map[string]any{
				"id":      c.GetId(),
				"uid":     c.GetUid(),
				"name":    profiles[c.GetUid()].name,
				"avatar":  profiles[c.GetUid()].avatar,
				"biz":     c.GetBiz(),
				"bizid":   c.GetBizid(),
				"content": c.GetContent(),
				"ctime":   c.GetCtime().AsTime().Format(time.DateTime),
				"utime":   c.GetUtime().AsTime().Format(time.DateTime),
			}
			if pc := c.GetParentComment(); pc != nil {
				child["parent_comment"] = map[string]any{
					"id":      pc.GetId(),
					"uid":     pc.GetUid(),
					"name":    profiles[pc.GetUid()].name,
					"avatar":  profiles[pc.GetUid()].avatar,
					"biz":     pc.GetBiz(),
					"bizid":   pc.GetBizid(),
					"content": pc.GetContent(),
					"ctime":   pc.GetCtime().AsTime().Format(time.DateTime),
					"utime":   pc.GetUtime().AsTime().Format(time.DateTime),
				}
			}
			root["children"] = append(root["children"].([]map[string]any), child)
		}
	}
	res := make([]map[string]any, 0, len(roots))
	for _, id := range order {
		if item := roots[id]; item != nil {
			res = append(res, item)
		}
	}
	return ginx.Result{Data: map[string]any{"comments": res}}, nil
}

func (h *CommentHandler) Create(ctx *gin.Context, req CreateReq, uc ginx.UserClaims) (ginx.Result, error) {
	comment := &commentv1.Comment{
		Uid:     uc.Id,
		Biz:     req.Comment.Biz,
		Bizid:   req.Comment.BizId,
		Content: req.Comment.Content,
	}
	if req.Comment.RootComment != nil {
		comment.RootComment = &commentv1.Comment{Id: req.Comment.RootComment.Id}
	}
	if req.Comment.ParentComment != nil {
		comment.ParentComment = &commentv1.Comment{Id: req.Comment.ParentComment.Id}
	}
	_, err := h.client.CreateComment(ctx, &commentv1.CreateCommentRequest{Comment: comment})
	if err != nil {
		return ginx.Result{Code: 5, Msg: "系统错误"}, err
	}
	return ginx.Result{Msg: "OK"}, nil
}

func (h *CommentHandler) Delete(ctx *gin.Context, req DeleteReq, uc ginx.UserClaims) (ginx.Result, error) {
	_, err := h.client.DeleteComment(ctx, &commentv1.DeleteCommentRequest{Id: req.Id})
	if err != nil {
		return ginx.Result{Code: 5, Msg: "系统错误"}, err
	}
	return ginx.Result{Msg: "OK"}, nil
}

func (h *CommentHandler) Replies(ctx *gin.Context, req RepliesReq) (ginx.Result, error) {
	resp, err := h.client.GetMoreReplies(ctx, &commentv1.GetMoreRepliesRequest{
		Rid:   req.Rid,
		MaxId: req.MaxId,
		Limit: req.Limit,
	})
	if err != nil {
		return ginx.Result{Code: 5, Msg: "系统错误"}, err
	}
	cs := resp.GetReplies()
	uids := make(map[int64]struct{}, len(cs))
	for _, c := range cs {
		if c.GetUid() > 0 {
			uids[c.GetUid()] = struct{}{}
		}
		if pc := c.GetParentComment(); pc != nil && pc.GetUid() > 0 {
			uids[pc.GetUid()] = struct{}{}
		}
		if rc := c.GetRootComment(); rc != nil && rc.GetUid() > 0 {
			uids[rc.GetUid()] = struct{}{}
		}
	}
	type ui struct{ name, avatar string }
	profiles := make(map[int64]ui, len(uids))
	for uid := range uids {
		pr, er := h.userSvc.Profile(ctx, &userv1.ProfileRequest{Id: uid})
		if er != nil || pr == nil || pr.User == nil {
			continue
		}
		profiles[uid] = ui{name: pr.User.Nickname, avatar: pr.User.Avatar}
	}
	res := make([]map[string]any, 0, len(cs))
	for _, c := range cs {
		item := map[string]any{
			"id":      c.GetId(),
			"uid":     c.GetUid(),
			"name":    profiles[c.GetUid()].name,
			"avatar":  profiles[c.GetUid()].avatar,
			"biz":     c.GetBiz(),
			"bizid":   c.GetBizid(),
			"content": c.GetContent(),
			"ctime":   c.GetCtime().AsTime().Format(time.DateTime),
			"utime":   c.GetUtime().AsTime().Format(time.DateTime),
		}
		if rc := c.GetRootComment(); rc != nil {
			item["root_comment"] = map[string]any{
				"id":      rc.GetId(),
				"uid":     rc.GetUid(),
				"name":    profiles[rc.GetUid()].name,
				"avatar":  profiles[rc.GetUid()].avatar,
				"biz":     rc.GetBiz(),
				"bizid":   rc.GetBizid(),
				"content": rc.GetContent(),
				"ctime":   rc.GetCtime().AsTime().Format(time.DateTime),
				"utime":   rc.GetUtime().AsTime().Format(time.DateTime),
			}
		}
		if pc := c.GetParentComment(); pc != nil {
			item["parent_comment"] = map[string]any{
				"id":      pc.GetId(),
				"uid":     pc.GetUid(),
				"name":    profiles[pc.GetUid()].name,
				"avatar":  profiles[pc.GetUid()].avatar,
				"biz":     pc.GetBiz(),
				"bizid":   pc.GetBizid(),
				"content": pc.GetContent(),
				"ctime":   pc.GetCtime().AsTime().Format(time.DateTime),
				"utime":   pc.GetUtime().AsTime().Format(time.DateTime),
			}
		}
		res = append(res, item)
	}
	return ginx.Result{Data: map[string]any{"replies": res}}, nil
}
