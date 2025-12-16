package web

import (
	intrv1 "Webook/api/proto/gen/api/proto/intr/v1"
	userv1 "Webook/api/proto/gen/api/proto/user/v1"
	"Webook/pkg/ginx"
	"Webook/pkg/logger"

	"github.com/gin-gonic/gin"
)

type NotificationHandler struct {
	l       logger.Logger
	cli     intrv1.InteractiveServiceClient
	userSvc userv1.UsersServiceClient
}

func NewNotificationHandler(l logger.Logger, cli intrv1.InteractiveServiceClient, userSvc userv1.UsersServiceClient) *NotificationHandler {
	return &NotificationHandler{l: l, cli: cli, userSvc: userSvc}
}

func (h *NotificationHandler) RegisterRoute(s *gin.Engine) {
	g := s.Group("/notifications")
	g.GET("/unread_counts", ginx.WrapClaims(h.UnreadCounts))
	g.GET("", ginx.WrapClaimsAndReq[NotificationListReq](h.List))
	g.POST("/mark_read", ginx.WrapClaimsAndReq[MarkReadReq](h.MarkRead))
}

type NotificationListReq struct {
	Type   string `json:"type" form:"type"`
	Offset int32  `json:"offset" form:"offset"`
	Limit  int32  `json:"limit" form:"limit"`
}

type MarkReadReq struct {
	Ids  []int64 `json:"ids"`
	Type string  `json:"type"`
}

func (h *NotificationHandler) UnreadCounts(ctx *gin.Context, uc ginx.UserClaims) (ginx.Result, error) {
	resp, err := h.cli.GetUnreadCounts(ctx, &intrv1.GetUnreadCountsRequest{Uid: uc.Id})
	if err != nil || resp == nil {
		return ginx.Result{Code: 5, Msg: "系统错误"}, err
	}
	return ginx.Result{Data: map[string]any{
		"interaction": resp.Interaction,
		"follow":      resp.Follow,
		"system":      resp.System,
		"total":       resp.Total,
	}}, nil
}

func (h *NotificationHandler) List(ctx *gin.Context, req NotificationListReq, uc ginx.UserClaims) (ginx.Result, error) {
	resp, err := h.cli.GetNotifications(ctx, &intrv1.GetNotificationsRequest{
		Uid:    uc.Id,
		Type:   req.Type,
		Offset: req.Offset,
		Limit:  req.Limit,
	})
	if err != nil || resp == nil {
		return ginx.Result{Code: 5, Msg: "系统错误"}, err
	}
	// 拉取发送者资料
	type ui struct{ name, avatar string }
	profiles := make(map[int64]ui)
	for _, it := range resp.Items {
		if it == nil || it.Sender == nil {
			continue
		}
		sid := it.Sender.Id
		if sid <= 0 {
			continue
		}
		if _, ok := profiles[sid]; ok {
			continue
		}
		pr, er := h.userSvc.Profile(ctx, &userv1.ProfileRequest{Id: sid})
		if er != nil || pr == nil || pr.User == nil {
			continue
		}
		u := pr.User
		profiles[sid] = ui{name: u.Nickname, avatar: u.Avatar}
	}

	items := make([]map[string]any, 0, len(resp.Items))
	for _, it := range resp.Items {
		if it == nil {
			continue
		}
		sender := map[string]any{}
		if it.Sender != nil {
			sid := it.Sender.Id
			sender["id"] = sid
			p := profiles[sid]
			sender["name"] = p.name
			sender["avatar"] = p.avatar
		} else {
			sender["id"] = int64(0)
			sender["name"] = ""
			sender["avatar"] = ""
		}
		var target map[string]any
		if it.Target != nil {
			target = map[string]any{
				"type":    it.Target.Type,
				"id":      it.Target.Id,
				"title":   it.Target.Title,
				"preview": it.Target.Preview,
			}
		}
		// 状态转数字：0 未读，1 已读
		status := 0
		if it.Status == "read" {
			status = 1
		}
		items = append(items, map[string]any{
			"id":       it.Id,
			"category": it.Category,
			"content":  it.Content,
			"time":     it.Time,
			"status":   status,
			"sender":   sender,
			"target":   target,
		})
	}
	return ginx.Result{Data: items}, nil
}

func (h *NotificationHandler) MarkRead(ctx *gin.Context, req MarkReadReq, uc ginx.UserClaims) (ginx.Result, error) {
	_, err := h.cli.MarkRead(ctx, &intrv1.MarkReadRequest{Uid: uc.Id, Ids: req.Ids, Type: req.Type})
	if err != nil {
		return ginx.Result{Code: 5, Msg: "系统错误"}, err
	}
	return ginx.Result{Msg: "OK"}, nil
}
