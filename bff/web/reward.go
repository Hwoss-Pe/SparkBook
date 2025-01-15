package web

import (
	articlev1 "Webook/api/proto/gen/api/proto/article/v1"
	rewardv1 "Webook/api/proto/gen/api/proto/reward/v1"
	"Webook/pkg/ginx"
	"github.com/gin-gonic/gin"
)

type RewardHandler struct {
	client    rewardv1.RewardServiceClient
	artClient articlev1.ArticleServiceClient
}

func (r *RewardHandler) RegisterRoute(server *gin.Engine) {
	rg := server.Group("/reward")
	rg.POST("/detail", ginx.WrapClaimsAndReq[GetRewardReq](r.GetReward))
}

func (r *RewardHandler) GetReward(ctx *gin.Context, req GetRewardReq, claims ginx.UserClaims) (ginx.Result, error) {
	resp, err := r.client.GetReward(ctx, &rewardv1.GetRewardRequest{
		Rid: req.Rid,
		Uid: claims.Id,
	})
	if err != nil {
		return ginx.Result{
			Code: 5,
			Msg:  "系统错误",
		}, err
	}
	return ginx.Result{
		Data: resp.Status.String(),
	}, nil
}

func NewRewardHandler(client rewardv1.RewardServiceClient, artClient articlev1.ArticleServiceClient) *RewardHandler {
	return &RewardHandler{client: client, artClient: artClient}
}

type GetRewardReq struct {
	Rid int64
}
type RewardArticleReq struct {
	Aid int64 `json:"aid"`
	Amt int64 `json:"amt"`
}
