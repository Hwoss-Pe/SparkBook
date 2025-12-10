package web

import (
	followv1 "Webook/api/proto/gen/api/proto/follow/v1"
	userv1 "Webook/api/proto/gen/api/proto/user/v1"
	"Webook/pkg/ginx"
	"Webook/pkg/logger"

	"github.com/gin-gonic/gin"
)

type FollowHandler struct {
	svc     followv1.FollowServiceClient
	userSvc userv1.UsersServiceClient
	l       logger.Logger
}

func NewFollowHandler(svc followv1.FollowServiceClient, userSvc userv1.UsersServiceClient, l logger.Logger) *FollowHandler {
	return &FollowHandler{svc: svc, userSvc: userSvc, l: l}
}

func (h *FollowHandler) RegisterRoute(s *gin.Engine) {
	g := s.Group("/follow")
	g.POST("", ginx.WrapClaimsAndReq[FollowReq](h.Follow))
	g.POST("/cancel", ginx.WrapClaimsAndReq[FollowReq](h.CancelFollow))
	g.GET("/followee", ginx.WrapReq[GetFolloweeReq](h.GetFollowee))
	g.GET("/follower", ginx.WrapReq[GetFollowerReq](h.GetFollower))
	g.GET("/info", ginx.WrapReq[FollowInfoReq](h.FollowInfo))
	g.GET("/statics", ginx.WrapReq[GetFollowStaticReq](h.GetFollowStatics))
}

type FollowReq struct {
	Followee int64 `json:"followee"`
	Follower int64 `json:"follower"`
}

type GetFolloweeReq struct {
	Follower int64 `form:"follower"`
	Offset   int64 `form:"offset"`
	Limit    int64 `form:"limit"`
}

type GetFollowerReq struct {
	Followee int64 `form:"followee"`
	Offset   int64 `form:"offset"`
	Limit    int64 `form:"limit"`
}

type FollowInfoReq struct {
	Follower int64 `form:"follower"`
	Followee int64 `form:"followee"`
}

type GetFollowStaticReq struct {
	Followee int64 `form:"followee"`
}

type FollowRelationVo struct {
	Id       int64  `json:"id"`
	Follower int64  `json:"follower"`
	Followee int64  `json:"followee"`
	Name     string `json:"name"`
	Avatar   string `json:"avatar"`
	AboutMe  string `json:"about_me"`
}

type FollowStaticVo struct {
	Followers int64 `json:"followers"`
	Followees int64 `json:"followees"`
}

func (h *FollowHandler) Follow(ctx *gin.Context, req FollowReq, uc ginx.UserClaims) (ginx.Result, error) {
	_, err := h.svc.Follow(ctx, &followv1.FollowRequest{Followee: req.Followee, Follower: uc.Id})
	if err != nil {
		return ginx.Result{Code: 5, Msg: "系统错误"}, err
	}
	return ginx.Result{Msg: "OK"}, nil
}

func (h *FollowHandler) CancelFollow(ctx *gin.Context, req FollowReq, uc ginx.UserClaims) (ginx.Result, error) {
	_, err := h.svc.CancelFollow(ctx, &followv1.CancelFollowRequest{Followee: req.Followee, Follower: uc.Id})
	if err != nil {
		return ginx.Result{Code: 5, Msg: "系统错误"}, err
	}
	return ginx.Result{Msg: "OK"}, nil
}

func (h *FollowHandler) GetFollowee(ctx *gin.Context, req GetFolloweeReq) (ginx.Result, error) {
	if req.Limit <= 0 || req.Limit > 100 {
		req.Limit = 100
	}
	resp, err := h.svc.GetFollowee(ctx, &followv1.GetFolloweeRequest{Follower: req.Follower, Offset: req.Offset, Limit: req.Limit})
	if err != nil {
		return ginx.Result{Code: 5, Msg: "系统错误"}, err
	}
	rels := resp.GetFollowRelations()
	// 收集 followee 的用户资料
	type ui struct{ name, avatar, aboutMe string }
	profiles := make(map[int64]ui, len(rels))
	for _, r := range rels {
		if r == nil {
			continue
		}
		uid := r.GetFollowee()
		if _, ok := profiles[uid]; ok {
			continue
		}
		pr, er := h.userSvc.Profile(ctx, &userv1.ProfileRequest{Id: uid})
		if er != nil || pr == nil || pr.User == nil {
			continue
		}
		u := pr.User
		profiles[uid] = ui{name: u.Nickname, avatar: u.Avatar, aboutMe: u.AboutMe}
	}
	data := make([]FollowRelationVo, 0, len(rels))
	for _, r := range rels {
		if r == nil {
			continue
		}
		p := profiles[r.GetFollowee()]
		data = append(data, FollowRelationVo{
			Id:       r.GetId(),
			Follower: r.GetFollower(),
			Followee: r.GetFollowee(),
			Name:     p.name,
			Avatar:   p.avatar,
			AboutMe:  p.aboutMe,
		})
	}
	return ginx.Result{Data: map[string]any{"follow_relations": data}}, nil
}

func (h *FollowHandler) GetFollower(ctx *gin.Context, req GetFollowerReq) (ginx.Result, error) {
	if req.Limit <= 0 || req.Limit > 100 {
		req.Limit = 100
	}
	resp, err := h.svc.GetFollower(ctx, &followv1.GetFollowerRequest{Followee: req.Followee, Offset: req.Offset, Limit: req.Limit})
	if err != nil {
		return ginx.Result{Code: 5, Msg: "系统错误"}, err
	}
	rels := resp.GetFollowRelations()
	// 收集 follower 的用户资料
	type ui struct{ name, avatar, aboutMe string }
	profiles := make(map[int64]ui, len(rels))
	for _, r := range rels {
		if r == nil {
			continue
		}
		uid := r.GetFollower()
		if _, ok := profiles[uid]; ok {
			continue
		}
		pr, er := h.userSvc.Profile(ctx, &userv1.ProfileRequest{Id: uid})
		if er != nil || pr == nil || pr.User == nil {
			continue
		}
		u := pr.User
		profiles[uid] = ui{name: u.Nickname, avatar: u.Avatar, aboutMe: u.AboutMe}
	}
	data := make([]FollowRelationVo, 0, len(rels))
	for _, r := range rels {
		if r == nil {
			continue
		}
		p := profiles[r.GetFollower()]
		data = append(data, FollowRelationVo{
			Id:       r.GetId(),
			Follower: r.GetFollower(),
			Followee: r.GetFollowee(),
			Name:     p.name,
			Avatar:   p.avatar,
			AboutMe:  p.aboutMe,
		})
	}
	return ginx.Result{Data: map[string]any{"follow_relations": data}}, nil
}

func (h *FollowHandler) FollowInfo(ctx *gin.Context, req FollowInfoReq) (ginx.Result, error) {
	resp, err := h.svc.FollowInfo(ctx, &followv1.FollowInfoRequest{Follower: req.Follower, Followee: req.Followee})
	if err != nil {
		return ginx.Result{Code: 5, Msg: "系统错误"}, err
	}
	if resp == nil || resp.FollowRelation == nil {
		return ginx.Result{Data: map[string]any{"follow_relation": nil}}, nil
	}
	fr := resp.FollowRelation
	return ginx.Result{Data: map[string]any{"follow_relation": FollowRelationVo{Id: fr.GetId(), Follower: fr.GetFollower(), Followee: fr.GetFollowee()}}}, nil
}

func (h *FollowHandler) GetFollowStatics(ctx *gin.Context, req GetFollowStaticReq) (ginx.Result, error) {
	resp, err := h.svc.GetFollowStatics(ctx, &followv1.GetFollowStaticRequest{Followee: req.Followee})
	if err != nil || resp == nil || resp.FollowStatic == nil {
		return ginx.Result{Code: 5, Msg: "系统错误"}, err
	}
	fs := resp.FollowStatic
	return ginx.Result{Data: map[string]any{"followStatic": FollowStaticVo{Followers: fs.GetFollowers(), Followees: fs.GetFollowees()}}}, nil
}
