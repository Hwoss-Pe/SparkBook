package web

import (
	articlev1 "Webook/api/proto/gen/api/proto/article/v1"
	intrv1 "Webook/api/proto/gen/api/proto/intr/v1"
	rewardv1 "Webook/api/proto/gen/api/proto/reward/v1"
	"Webook/bff/web/jwt"
	"Webook/pkg/ginx"
	"Webook/pkg/logger"
	"fmt"
	"github.com/ecodeclub/ekit/slice"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
	"google.golang.org/protobuf/types/known/timestamppb"
	"net/http"
	"strconv"
	"time"
)

type ArticleHandler struct {
	svc     articlev1.ArticleServiceClient
	intrSvc intrv1.InteractiveServiceClient
	reward  rewardv1.RewardServiceClient
	l       logger.Logger
	biz     string
}

func NewArticleHandler(svc articlev1.ArticleServiceClient,
	intrSvc intrv1.InteractiveServiceClient,
	reward rewardv1.RewardServiceClient,
	l logger.Logger) *ArticleHandler {
	return &ArticleHandler{
		svc:     svc,
		l:       l,
		reward:  reward,
		biz:     "article",
		intrSvc: intrSvc,
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
	// 推荐文章列表（首页使用）
	pub.POST("/list", ginx.WrapReq[ListPubReq](a.ListPub))
	pub.POST("/like", ginx.WrapClaimsAndReq[LikeReq](a.Like))
	pub.POST("/collect", ginx.WrapClaimsAndReq[CollectReq](a.Collect))
	// 打赏
	pub.POST("/reward", ginx.WrapClaimsAndReq[RewardReq](a.Reward))
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

	go func() {
		_, err = a.intrSvc.IncrReadCnt(ctx, &intrv1.IncrReadCntRequest{
			Biz:   "article",
			BizId: art.Id,
		})
		if err != nil {
			a.l.Error("增加文章阅读数失败", logger.Error(err))
		}
	}()
	intr := intrResp.Intr
	return Result{
		Data: ArticleVo{
			Id:         art.Id,
			Title:      art.Title,
			Status:     art.Status,
			Content:    art.Content,
			CoverImage: art.CoverImage,
			// 要把作者信息带出去
			Author:     art.Author.Name,
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
		// 设置作者信息
		if art.Author != nil {
			vo.Author = AuthorVo{
				Id:   art.Author.Id,
				Name: art.Author.Name,
			}
		}
		// 设置互动数据
		if intr, ok := intrMap[art.Id]; ok {
			vo.ReadCnt = intr.ReadCnt
			vo.LikeCnt = intr.LikeCnt
			vo.CollectCnt = intr.CollectCnt
		}
		result = append(result, vo)
	}

	return ginx.Result{
		Data: result,
	}, nil
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
