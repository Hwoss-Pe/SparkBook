package web

import (
	articlev1 "Webook/api/proto/gen/api/proto/article/v1"
	intrv1 "Webook/api/proto/gen/api/proto/intr/v1"
	rankingv1 "Webook/api/proto/gen/api/proto/ranking/v1"
	rewardv1 "Webook/api/proto/gen/api/proto/reward/v1"
	"Webook/bff/web/jwt"
	"Webook/pkg/ginx"
	"Webook/pkg/logger"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/ecodeclub/ekit/slice"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ArticleHandler struct {
	svc        articlev1.ArticleServiceClient
	intrSvc    intrv1.InteractiveServiceClient
	rankingSvc rankingv1.RankingServiceClient
	reward     rewardv1.RewardServiceClient
	l          logger.Logger
	biz        string
}

func NewArticleHandler(svc articlev1.ArticleServiceClient,
	intrSvc intrv1.InteractiveServiceClient,
	rankingSvc rankingv1.RankingServiceClient,
	reward rewardv1.RewardServiceClient,
	l logger.Logger) *ArticleHandler {
	return &ArticleHandler{
		svc:        svc,
		l:          l,
		reward:     reward,
		biz:        "article",
		intrSvc:    intrSvc,
		rankingSvc: rankingSvc,
	}
}

func (a *ArticleHandler) RegisterRoute(s *gin.Engine) {
	g := s.Group("/articles")
	g.GET("/detail/:id", a.Detail)
	g.POST("/list", ginx.WrapClaimsAndReq(a.List))
	g.POST("/edit", a.Edit)
	g.POST("/publish", a.Publish)
	g.POST("/withdraw", a.Withdraw)
	pub := g.Group("/pub")
	pub.GET("/:id", ginx.WrapClaims(a.PubDetail))
	// 推荐文章列表（首页使用，匿名可访问）
	pub.POST("/list", ginx.WrapReq[ListPubReq](a.ListPub))
	pub.POST("/like", ginx.WrapClaimsAndReq[LikeReq](a.Like))
	pub.POST("/cancelLike", ginx.WrapClaimsAndReq[LikeReq](a.CancelLike))
	pub.POST("/collect", ginx.WrapClaimsAndReq[CollectReq](a.Collect))
	pub.POST("/cancelCollect", ginx.WrapClaimsAndReq[CollectReq](a.CancelCollect))
	// 打赏
	pub.POST("/reward", ginx.WrapClaimsAndReq[RewardReq](a.Reward))
	// 热榜（匿名可访问）
	pub.GET("/ranking", ginx.WrapReq[RankingReq](a.Ranking))
	// 手动触发热榜计算
	pub.POST("/ranking/trigger", ginx.WrapReq[TriggerRankingReq](a.TriggerRanking))
}

func (a *ArticleHandler) Detail(ctx *gin.Context) {
	idstr := ctx.Param("id")
	id, err := strconv.ParseInt(idstr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "参数错误",
		})
		a.l.Error("前端输入的 ID 不对", logger.Error(err))
		return
	}
	usr, ok := ctx.MustGet("user").(jwt.UserClaims)
	if !ok {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		a.l.Error("获得用户会话信息失败")
		return
	}
	resp, err := a.svc.GetById(ctx, &articlev1.GetByIdRequest{Id: id})
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		a.l.Error("获得文章信息失败", logger.Error(err))
		return
	}
	art := resp.GetArticle()

	//保证文章只能看到自己写的
	if art.Author.Id != usr.Id {
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			// 也不需要告诉前端究竟发生了什么
			Msg: "输入有误",
		})
		// 如果公司有风控系统，这个时候就要上报这种非法访问的用户了。
		a.l.Error("非法访问文章，创作者 ID 不匹配",
			logger.Int64("uid", usr.Id))
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Data: ArticleVo{
			Id:         art.Id,
			Title:      art.Title,
			Status:     art.Status,
			Content:    art.Content,
			CoverImage: art.CoverImage,
			Ctime:      art.Ctime.AsTime().Format(time.DateTime),
			Utime:      art.Utime.AsTime().Format(time.DateTime),
		},
	})
}

type ListReq struct {
	Offset int32 `json:"offset"`
	Limit  int32 `json:"limit"`
}

// ListPubReq 推荐文章列表请求
type ListPubReq struct {
	Offset int32 `json:"offset"`
	Limit  int32 `json:"limit"`
}

func (a *ArticleHandler) List(ctx *gin.Context, req ListReq, usr ginx.UserClaims) (ginx.Result, error) {
	if req.Limit > 100 {
		a.l.Error("获得用户会话信息失败，LIMIT过大")
		return ginx.Result{
			Code: 4,
			Msg:  "请求有误",
		}, nil
	}
	arts, err := a.svc.List(ctx, &articlev1.ListRequest{Author: usr.Id,
		Offset: req.Offset, Limit: req.Limit})
	if err != nil {
		a.l.Error("获得用户会话信息失败")
		return ginx.Result{
			Code: 5,
			Msg:  "系统错误",
		}, nil
	}
	return ginx.Result{
		Data: slice.Map[*articlev1.Article, ArticleVo](arts.Articles,
			func(idx int, src *articlev1.Article) ArticleVo {
				return ArticleVo{
					Id:         src.Id,
					Title:      src.Title,
					Abstract:   src.Abstract,
					CoverImage: src.CoverImage,
					Status:     src.Status,
					Ctime:      src.Ctime.AsTime().Format(time.DateTime),
					Utime:      src.Utime.AsTime().Format(time.DateTime),
				}
			}),
	}, nil
}

func (a *ArticleHandler) Edit(ctx *gin.Context) {
	var req ArticleReq
	if err := ctx.Bind(&req); err != nil {
		a.l.Error("反序列化请求失败", logger.Error(err))
		return
	}
	usr, ok := ctx.MustGet("user").(jwt.UserClaims)
	if !ok {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		a.l.Error("获得用户会话信息失败")
		return
	}
	id, err := a.svc.Save(ctx, &articlev1.SaveRequest{Article: req.toDTO(usr.Id)})
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		a.l.Error("保存数据失败", logger.Field{Key: "error", Value: err})
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Data: id,
	})
}

func (a *ArticleHandler) Publish(ctx *gin.Context) {
	var req ArticleReq
	if err := ctx.Bind(&req); err != nil {
		a.l.Error("反序列化请求失败", logger.Error(err))
		return
	}
	usr, ok := ctx.MustGet("user").(jwt.UserClaims)
	if !ok {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		a.l.Error("获得用户会话信息失败")
		return
	}
	idResp, err := a.svc.Publish(ctx, &articlev1.PublishRequest{
		Article: &articlev1.Article{
			Id:         req.Id,
			Title:      req.Title,
			Content:    req.Content,
			CoverImage: req.CoverImage,
			Author: &articlev1.Author{
				Id: usr.Id,
			},
		},
	})
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		a.l.Error("发表失败", logger.Error(err))
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Data: idResp.Id,
	})
}

func (a *ArticleHandler) PubDetail(ctx *gin.Context, uc ginx.UserClaims) (ginx.Result, error) {
	idstr := ctx.Param("id")
	id, err := strconv.ParseInt(idstr, 10, 64)
	if err != nil {
		a.l.Error("前端输入的 ID 不对", logger.Error(err))
		return Result{
			Code: 4,
			Msg:  "参数错误",
		}, fmt.Errorf("查询文章详情的 ID %s 不正确, %w", idstr, err)
	}
	var (
		eg       errgroup.Group
		artResp  *articlev1.GetPublishedByIdResponse
		intrResp *intrv1.GetResponse
	)
	eg.Go(func() error {
		var er error
		artResp, er = a.svc.GetPublishedById(ctx, &articlev1.GetPublishedByIdRequest{
			Id: id, Uid: uc.Id,
		})
		return er
	})

	eg.Go(func() error {
		var er error
		intrResp, er = a.intrSvc.Get(ctx, &intrv1.GetRequest{
			Biz: a.biz, BizId: id, Uid: uc.Id,
		})
		return er
	})

	if err != nil {
		return Result{
			Code: 5,
			Msg:  "系统错误",
		}, fmt.Errorf("获取文章信息失败 %w", err)
	}
	err = eg.Wait()
	if err != nil {
		return ginx.Result{}, err
	}
	art := artResp.GetArticle()
	if art == nil {
		return Result{Data: nil}, nil
	}

	// 阅读量不在BFF直接递增，统一走事件驱动链路
	intr := intrResp.Intr
	return Result{
		Data: ArticleVo{
			Id:         art.Id,
			Title:      art.Title,
			Status:     art.Status,
			Content:    art.Content,
			CoverImage: art.CoverImage,
			Author:     AuthorVo{Id: art.Author.GetId(), Name: art.Author.GetName(), Avatar: art.Author.GetAvatar()},
			Ctime:      art.Ctime.AsTime().Format(time.DateTime),
			Utime:      art.Utime.AsTime().Format(time.DateTime),
			ReadCnt:    intr.ReadCnt,
			CollectCnt: intr.CollectCnt,
			LikeCnt:    intr.LikeCnt,
			Liked:      intr.Liked,
			Collected:  intr.Collected,
		},
	}, nil
}

func (a *ArticleHandler) Withdraw(ctx *gin.Context) {

}

// ListPub 获取推荐文章列表（首页使用，按时间排序）
func (a *ArticleHandler) ListPub(ctx *gin.Context, req ListPubReq) (ginx.Result, error) {
	if req.Limit > 100 || req.Limit <= 0 {
		req.Limit = 10
	}
	// 获取已发布文章列表（按时间排序）
	artResp, err := a.svc.ListPub(ctx, &articlev1.ListPubRequest{
		StartTime: timestamppb.Now(),
		Offset:    req.Offset,
		Limit:     req.Limit,
	})
	if err != nil {
		a.l.Error("获取推荐文章列表失败", logger.Error(err))
		return ginx.Result{
			Code: 5,
			Msg:  "系统错误",
		}, err
	}

	articles := artResp.GetArticles()
	if len(articles) == 0 {
		return ginx.Result{
			Data: []ArticlePubVo{},
		}, nil
	}

	// 提取文章ID列表
	ids := make([]int64, 0, len(articles))
	for _, art := range articles {
		ids = append(ids, art.Id)
	}

	// 批量获取互动数据
	intrResp, err := a.intrSvc.GetByIds(ctx, &intrv1.GetByIdsRequest{
		Biz: a.biz,
		Ids: ids,
	})
	if err != nil {
		a.l.Error("获取互动数据失败", logger.Error(err))
		// 互动数据获取失败不影响文章列表返回，只是互动数据为空
	}

	intrMap := make(map[int64]*intrv1.Interactive)
	if intrResp != nil {
		intrMap = intrResp.GetIntrs()
	}

	type userIntr struct{ liked, collected bool }
	userMap := make(map[int64]userIntr, len(ids))
	var uid int64
	if raw, ok := ctx.Get("user"); ok {
		if claims, ok2 := raw.(ginx.UserClaims); ok2 {
			uid = claims.Id
		}
	}
	if uid > 0 {
		for _, id := range ids {
			resp, er := a.intrSvc.Get(ctx, &intrv1.GetRequest{Biz: a.biz, BizId: id, Uid: uid})
			if er != nil || resp == nil || resp.Intr == nil {
				continue
			}
			userMap[id] = userIntr{liked: resp.Intr.Liked, collected: resp.Intr.Collected}
		}
	}

	// 组装返回数据
	result := make([]ArticlePubVo, 0, len(articles))
	for _, art := range articles {
		vo := ArticlePubVo{
			Id:         art.Id,
			Title:      art.Title,
			Abstract:   art.Abstract,
			CoverImage: art.CoverImage,
			Ctime:      art.Ctime.AsTime().Format(time.DateTime),
			Utime:      art.Utime.AsTime().Format(time.DateTime),
		}
		if art.Author != nil {
			vo.Author = AuthorVo{Id: art.Author.Id, Name: art.Author.Name, Avatar: art.Author.Avatar}
		}
		if intr, ok := intrMap[art.Id]; ok {
			vo.ReadCnt = intr.ReadCnt
			vo.LikeCnt = intr.LikeCnt
			vo.CollectCnt = intr.CollectCnt
		}
		if u, ok := userMap[art.Id]; ok {
			vo.Liked = u.liked
			vo.Collected = u.collected
		}
		result = append(result, vo)
	}

	return ginx.Result{
		Data: result,
	}, nil
}

func (a *ArticleHandler) CancelLike(ctx *gin.Context, req LikeReq, uc ginx.UserClaims) (ginx.Result, error) {
	_, err := a.intrSvc.CancelLike(ctx, &intrv1.CancelLikeRequest{Biz: a.biz, BizId: req.Id, Uid: uc.Id})
	if err != nil {
		return Result{Code: 5, Msg: "系统错误"}, err
	}
	return Result{Msg: "OK"}, nil
}

func (a *ArticleHandler) Like(ctx *gin.Context, req LikeReq, uc ginx.UserClaims) (ginx.Result, error) {
	var err error
	if req.Like {
		_, err = a.intrSvc.Like(ctx, &intrv1.LikeRequest{
			Biz: a.biz, BizId: req.Id, Uid: uc.Id,
		})
	} else {
		_, err = a.intrSvc.CancelLike(ctx, &intrv1.CancelLikeRequest{
			Biz: a.biz, BizId: req.Id, Uid: uc.Id,
		})
	}
	if err != nil {
		return Result{
			Code: 5,
			Msg:  "系统错误",
		}, err
	}
	return Result{Msg: "OK"}, nil
}

func (a *ArticleHandler) Collect(ctx *gin.Context, req CollectReq, uc ginx.UserClaims) (ginx.Result, error) {
	_, err := a.intrSvc.Collect(ctx, &intrv1.CollectRequest{
		Biz: a.biz, BizId: req.Id, Uid: uc.Id,
		Cid: req.Cid,
	})
	if err != nil {
		return Result{
			Code: 5,
			Msg:  "系统错误",
		}, err
	}
	return Result{Msg: "OK"}, nil
}

func (a *ArticleHandler) CancelCollect(ctx *gin.Context, req CollectReq, uc ginx.UserClaims) (ginx.Result, error) {
	_, err := a.intrSvc.CancelCollect(ctx, &intrv1.CancelCollectRequest{
		Biz: a.biz, BizId: req.Id, Uid: uc.Id, Cid: req.Cid,
	})
	if err != nil {
		return Result{Code: 5, Msg: "系统错误"}, err
	}
	return Result{Msg: "OK"}, nil
}

func (a *ArticleHandler) Reward(ctx *gin.Context, req RewardReq, uc ginx.UserClaims) (ginx.Result, error) {
	artResp, err := a.svc.GetPublishedById(ctx.Request.Context(), &articlev1.GetPublishedByIdRequest{
		Id: req.Id,
	})
	if err != nil {
		return ginx.Result{
			Code: 5,
			Msg:  "系统错误",
		}, err
	}
	art := artResp.GetArticle()

	resp, err := a.reward.PreReward(ctx.Request.Context(), &rewardv1.PreRewardRequest{
		Biz:       "article",
		BizId:     art.Id,
		BizName:   art.Title,
		TargetUid: art.Author.GetId(),
		Uid:       uc.Id,
		Amt:       req.Amt,
	})
	if err != nil {
		return ginx.Result{
			Code: 5,
			Msg:  "系统错误",
		}, err
	}
	return ginx.Result{
		Data: map[string]any{
			"codeURL": resp.CodeUrl,
			"rid":     resp.Rid,
		},
	}, nil
}

// RankingReq 热榜请求
type RankingReq struct {
	Offset int32 `json:"offset"`
	Limit  int32 `json:"limit"`
}

// TriggerRankingReq 触发热榜计算请求
type TriggerRankingReq struct {
}

// TriggerRanking 手动触发热榜计算
func (a *ArticleHandler) TriggerRanking(ctx *gin.Context, req TriggerRankingReq) (ginx.Result, error) {
	// 调用 ranking 服务的 RankTopN 方法来触发热榜计算
	_, err := a.rankingSvc.RankTopN(ctx, &rankingv1.RankTopNRequest{})
	if err != nil {
		a.l.Error("触发热榜计算失败", logger.Error(err))
		return ginx.Result{
			Code: 5,
			Msg:  "触发热榜计算失败",
		}, err
	}

	a.l.Info("热榜计算已手动触发")
	return ginx.Result{
		Code: 0,
		Msg:  "热榜计算已成功触发",
		Data: map[string]any{
			"message": "热榜正在重新计算中，请稍后查看最新结果",
		},
	}, nil
}

// Ranking 获取热榜文章
func (a *ArticleHandler) Ranking(ctx *gin.Context, req RankingReq) (ginx.Result, error) {
	if req.Limit > 100 || req.Limit <= 0 {
		req.Limit = 10
	}

	// 调用热榜服务获取热榜文章
	rankResp, err := a.rankingSvc.TopN(ctx, &rankingv1.TopNRequest{})
	if err != nil {
		a.l.Error("获取热榜失败", logger.Error(err))
		return ginx.Result{
			Code: 5,
			Msg:  "系统错误",
		}, err
	}

	allArticles := rankResp.GetArticles()
	if len(allArticles) == 0 {
		return ginx.Result{
			Data: []ArticlePubVo{},
		}, nil
	}

	// 在BFF层进行分页处理
	start := int(req.Offset)
	end := start + int(req.Limit)
	if start >= len(allArticles) {
		return ginx.Result{
			Data: []ArticlePubVo{},
		}, nil
	}
	if end > len(allArticles) {
		end = len(allArticles)
	}
	articles := allArticles[start:end]

	// 提取文章ID列表
	ids := make([]int64, 0, len(articles))
	for _, art := range articles {
		ids = append(ids, art.Id)
	}

	// 批量获取互动数据
	intrResp, err := a.intrSvc.GetByIds(ctx, &intrv1.GetByIdsRequest{
		Biz: a.biz,
		Ids: ids,
	})
	if err != nil {
		a.l.Error("获取互动数据失败", logger.Error(err))
		// 互动数据获取失败不影响文章列表返回，只是互动数据为空
	}

	intrMap := make(map[int64]*intrv1.Interactive)
	if intrResp != nil {
		intrMap = intrResp.GetIntrs()
	}

	type userIntr struct{ liked, collected bool }
	userMap := make(map[int64]userIntr, len(ids))
	var uid int64
	if raw, ok := ctx.Get("user"); ok {
		if claims, ok2 := raw.(ginx.UserClaims); ok2 {
			uid = claims.Id
		}
	}
	if uid > 0 {
		for _, id := range ids {
			resp, er := a.intrSvc.Get(ctx, &intrv1.GetRequest{Biz: a.biz, BizId: id, Uid: uid})
			if er != nil || resp == nil || resp.Intr == nil {
				continue
			}
			userMap[id] = userIntr{liked: resp.Intr.Liked, collected: resp.Intr.Collected}
		}
	}

	// 组装返回数据
	result := make([]ArticlePubVo, 0, len(articles))
	for _, art := range articles {
		vo := ArticlePubVo{
			Id:         art.Id,
			Title:      art.Title,
			CoverImage: art.CoverImage,
			Ctime:      art.Ctime.AsTime().Format(time.DateTime),
			Utime:      art.Utime.AsTime().Format(time.DateTime),
		}
		if art.Author != nil {
			vo.Author = AuthorVo{Id: art.Author.Id, Name: art.Author.Name}
		}
		if intr, ok := intrMap[art.Id]; ok {
			vo.ReadCnt = intr.ReadCnt
			vo.LikeCnt = intr.LikeCnt
			vo.CollectCnt = intr.CollectCnt
		}
		if u, ok := userMap[art.Id]; ok {
			vo.Liked = u.liked
			vo.Collected = u.collected
		}
		result = append(result, vo)
	}

	return ginx.Result{
		Data: result,
	}, nil
}
