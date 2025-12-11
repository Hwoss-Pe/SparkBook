package web

import (
	articlev1 "Webook/api/proto/gen/api/proto/article/v1"
	followv1 "Webook/api/proto/gen/api/proto/follow/v1"
	intrv1 "Webook/api/proto/gen/api/proto/intr/v1"
	rankingv1 "Webook/api/proto/gen/api/proto/ranking/v1"
	rewardv1 "Webook/api/proto/gen/api/proto/reward/v1"
	searchv1 "Webook/api/proto/gen/api/proto/search/v1"
	tagv1 "Webook/api/proto/gen/api/proto/tag/v1"
	"Webook/bff/web/jwt"
	"Webook/pkg/ginx"
	"Webook/pkg/logger"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/ecodeclub/ekit/slice"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/resolver"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ArticleHandler struct {
	svc        articlev1.ArticleServiceClient
	intrSvc    intrv1.InteractiveServiceClient
	followSvc  followv1.FollowServiceClient
	rankingSvc rankingv1.RankingServiceClient
	reward     rewardv1.RewardServiceClient
	l          logger.Logger
	biz        string
	tagSvc     tagv1.TagServiceClient
	searchSvc  searchv1.SearchServiceClient
}

func NewArticleHandler(svc articlev1.ArticleServiceClient,
	intrSvc intrv1.InteractiveServiceClient,
	rankingSvc rankingv1.RankingServiceClient,
	reward rewardv1.RewardServiceClient,
	l logger.Logger,
	followSvc followv1.FollowServiceClient,
	tagSvc tagv1.TagServiceClient) *ArticleHandler {
	return &ArticleHandler{
		svc:        svc,
		l:          l,
		reward:     reward,
		biz:        "article",
		intrSvc:    intrSvc,
		rankingSvc: rankingSvc,
		followSvc:  followSvc,
		tagSvc:     tagSvc,
	}
}

func (a *ArticleHandler) RegisterRoute(s *gin.Engine) {
	g := s.Group("/articles")
	g.GET("/detail/:id", a.Detail)
	g.GET("/tags/official", a.OfficialTags)
	g.POST("/list", ginx.WrapClaimsAndReq(a.List))
	g.POST("/edit", a.Edit)
	g.POST("/publish", a.Publish)
	g.POST("/withdraw", a.Withdraw)
	g.POST("/unpublish", a.Unpublish)
	g.GET("/author/:id/stats", a.AuthorStats)
	pub := g.Group("/pub")
	pub.GET("/:id", ginx.WrapClaims(a.PubDetail))
	// 推荐文章列表（首页使用，匿名可访问）
	pub.POST("/list", ginx.WrapReq[ListPubReq](a.ListPub))
	// 作者的已发布文章列表（匿名可访问）
	pub.GET("/author/:id/list", ginx.WrapReq[AuthorPubListReq](a.ListPubByAuthor))
	// 获取用户收藏的文章列表（需要登录）
	pub.GET("/collected/list", ginx.WrapClaimsAndReq[CollectedListReq](a.ListCollected))
	pub.GET("/following/list", ginx.WrapClaimsAndReq[FollowingListReq](a.ListFollowing))
	pub.POST("/like", ginx.WrapClaimsAndReq[LikeReq](a.Like))
	pub.POST("/cancelLike", ginx.WrapClaimsAndReq[LikeReq](a.CancelLike))
	pub.POST("/collect", ginx.WrapClaimsAndReq[CollectReq](a.Collect))
	pub.POST("/cancelCollect", ginx.WrapClaimsAndReq[CollectReq](a.CancelCollect))
	// 打赏
	pub.POST("/reward", ginx.WrapClaimsAndReq[RewardReq](a.Reward))
	// 热榜（匿名可访问）
	pub.GET("/ranking", ginx.WrapReq[RankingReq](a.Ranking))
	pub.GET("/ranking/tag", ginx.WrapReq[TagRankingReq](a.RankingByTag))
	// 根据官方标签获取文章列表（匿名可访问）
	pub.GET("/tag/articles", ginx.WrapReq[OfficialTagListReq](a.ListByOfficialTag))
	// 手动触发热榜计算
	pub.POST("/ranking/trigger", ginx.WrapReq[TriggerRankingReq](a.TriggerRanking))
}

func (a *ArticleHandler) AuthorStats(ctx *gin.Context) {
	idstr := ctx.Param("id")
	id, err := strconv.ParseInt(idstr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{Code: 4, Msg: "参数错误"})
		return
	}
	var publishedIds []int64
	var publishedCount int64
	var draftCount int64
	const pageSize int32 = 100
	for offset := int32(0); ; offset += pageSize {
		resp, er := a.svc.List(ctx, &articlev1.ListRequest{Author: id, Offset: offset, Limit: pageSize})
		if er != nil {
			ctx.JSON(http.StatusOK, Result{Code: 5, Msg: "系统错误"})
			return
		}
		arts := resp.GetArticles()
		if len(arts) == 0 {
			break
		}
		for _, art := range arts {
			if art.Status == 2 || art.Status == 3 {
				publishedCount++
				publishedIds = append(publishedIds, art.Id)
			} else {
				draftCount++
			}
		}
		if len(arts) < int(pageSize) {
			break
		}
	}
	var totalRead int64
	var totalLike int64
	for i := 0; i < len(publishedIds); i += 50 {
		end := i + 50
		if end > len(publishedIds) {
			end = len(publishedIds)
		}
		batch := publishedIds[i:end]
		intrResp, er := a.intrSvc.GetByIds(ctx, &intrv1.GetByIdsRequest{Biz: a.biz, Ids: batch})
		if er != nil || intrResp == nil {
			continue
		}
		for _, intr := range intrResp.GetIntrs() {
			if intr != nil {
				totalRead += intr.ReadCnt
				totalLike += intr.LikeCnt
			}
		}
	}
	var followingCount int64
	var followerCount int64
	if a.followSvc != nil {
		fsResp, er := a.followSvc.GetFollowStatics(ctx, &followv1.GetFollowStaticRequest{Followee: id})
		if er == nil && fsResp != nil && fsResp.FollowStatic != nil {
			followerCount = fsResp.FollowStatic.Followers
			followingCount = fsResp.FollowStatic.Followees
		}
	}
	ctx.JSON(http.StatusOK, Result{Data: AuthorStatsVo{
		PublishedCount: publishedCount,
		DraftCount:     draftCount,
		TotalReadCount: totalRead,
		TotalLikeCount: totalLike,
		FollowingCount: followingCount,
		FollowerCount:  followerCount,
	}})
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
	var tags []string
	if a.tagSvc != nil && art.Author != nil {
		tgResp, _ := a.tagSvc.GetBizTags(ctx, &tagv1.GetBizTagsRequest{Uid: usr.Id, Biz: a.biz, BizId: id})
		if tgResp != nil {
			for _, t := range tgResp.Tags {
				if t != nil {
					tags = append(tags, t.Name)
				}
			}
		}
	}
	ctx.JSON(http.StatusOK, Result{
		Data: ArticleVo{
			Id:         art.Id,
			Title:      art.Title,
			Status:     art.Status,
			Content:    art.Content,
			CoverImage: art.CoverImage,
			Ctime:      formatProtoTime(art.Ctime),
			Utime:      formatProtoTime(art.Utime),
			Tags:       tags,
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

// AuthorPubListReq 作者已发布文章列表请求
type AuthorPubListReq struct {
	Offset int32 `json:"offset" form:"offset"`
	Limit  int32 `json:"limit" form:"limit"`
}

// CollectedListReq 用户收藏文章列表请求
type CollectedListReq struct {
	Offset int32 `json:"offset" form:"offset"`
	Limit  int32 `json:"limit" form:"limit"`
}

type FollowingListReq struct {
	Offset int32 `json:"offset" form:"offset"`
	Limit  int32 `json:"limit" form:"limit"`
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
	// 批量获取互动数据（阅读、点赞、收藏）
	ids := make([]int64, 0, len(arts.Articles))
	for _, art := range arts.Articles {
		ids = append(ids, art.Id)
	}
	intrMap := make(map[int64]*intrv1.Interactive)
	if len(ids) > 0 {
		intrResp, er := a.intrSvc.GetByIds(ctx, &intrv1.GetByIdsRequest{Biz: a.biz, Ids: ids})
		if er != nil {
			a.l.Error("获取互动数据失败", logger.Error(er))
		} else if intrResp != nil {
			intrMap = intrResp.GetIntrs()
		}
	}

	// 组装返回，合并互动统计
	data := slice.Map[*articlev1.Article, ArticleVo](arts.Articles,
		func(idx int, src *articlev1.Article) ArticleVo {
			vo := ArticleVo{
				Id:         src.Id,
				Title:      src.Title,
				Abstract:   src.Abstract,
				CoverImage: src.CoverImage,
				Status:     src.Status,
				Ctime:      formatProtoTime(src.Ctime),
				Utime:      formatProtoTime(src.Utime),
			}
			if intr, ok := intrMap[src.Id]; ok && intr != nil {
				vo.ReadCnt = intr.ReadCnt
				vo.LikeCnt = intr.LikeCnt
				vo.CollectCnt = intr.CollectCnt
			}
			return vo
		})
	return ginx.Result{Data: data}, nil
}

func formatProtoTime(ts *timestamppb.Timestamp) string {
	if ts == nil {
		return "-"
	}
	t := ts.AsTime()
	if t.IsZero() || t.Year() < 1971 {
		return "-"
	}
	return t.Format(time.DateTime)
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
	if a.tagSvc != nil && len(req.Tags) > 0 {
		names := make(map[string]struct{}, len(req.Tags))
		for _, n := range req.Tags {
			if n == "" {
				continue
			}
			names[n] = struct{}{}
		}
		var existed map[string]int64
		existed = make(map[string]int64)
		resp, er := a.tagSvc.GetTags(ctx, &tagv1.GetTagsRequest{Uid: usr.Id})
		if er == nil && resp != nil {
			for _, t := range resp.Tag {
				existed[t.Name] = t.Id
			}
		}
		ids := make([]int64, 0, len(names))
		for n := range names {
			if id, ok := existed[n]; ok {
				ids = append(ids, id)
				continue
			}
			cr, e2 := a.tagSvc.CreateTag(ctx, &tagv1.CreateTagRequest{Uid: usr.Id, Name: n})
			if e2 != nil || cr == nil || cr.Tag == nil {
				continue
			}
			ids = append(ids, cr.Tag.Id)
		}
		if len(ids) > 0 {
			_, _ = a.tagSvc.AttachTags(ctx, &tagv1.AttachTagsRequest{Uid: usr.Id, Biz: a.biz, BizId: idResp.Id, Tids: ids})
		}
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
	var tags []string
	if a.tagSvc != nil && art.Author != nil {
		tgResp, _ := a.tagSvc.GetBizTags(ctx, &tagv1.GetBizTagsRequest{Uid: art.Author.GetId(), Biz: a.biz, BizId: id})
		if tgResp != nil {
			for _, t := range tgResp.Tags {
				if t != nil {
					tags = append(tags, t.Name)
				}
			}
		}
	}
	return Result{
		Data: ArticleVo{
			Id:         art.Id,
			Title:      art.Title,
			Status:     art.Status,
			Content:    art.Content,
			CoverImage: art.CoverImage,
			Author:     AuthorVo{Id: art.Author.GetId(), Name: art.Author.GetName(), Avatar: art.Author.GetAvatar()},
			Ctime:      formatProtoTime(art.Ctime),
			Utime:      formatProtoTime(art.Utime),
			ReadCnt:    intr.ReadCnt,
			CollectCnt: intr.CollectCnt,
			LikeCnt:    intr.LikeCnt,
			Liked:      intr.Liked,
			Collected:  intr.Collected,
			Tags:       tags,
		},
	}, nil
}

func (a *ArticleHandler) Withdraw(ctx *gin.Context) {
	var req struct {
		Id  int64 `json:"id"`
		Uid int64 `json:"uid"`
	}
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
	_, err := a.svc.Withdraw(ctx, &articlev1.WithdrawRequest{Id: req.Id, Uid: usr.Id})
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		a.l.Error("撤回失败", logger.Error(err))
		return
	}
	ctx.JSON(http.StatusOK, Result{Msg: "OK"})
}

func (a *ArticleHandler) Unpublish(ctx *gin.Context) {
	var req struct {
		Id  int64 `json:"id"`
		Uid int64 `json:"uid"`
	}
	if err := ctx.Bind(&req); err != nil {
		a.l.Error("反序列化请求失败", logger.Error(err))
		return
	}
	usr, ok := ctx.MustGet("user").(jwt.UserClaims)
	if !ok {
		ctx.JSON(http.StatusOK, Result{Code: 5, Msg: "系统错误"})
		a.l.Error("获得用户会话信息失败")
		return
	}
	_, err := a.svc.Unpublish(ctx, &articlev1.UnpublishRequest{Id: req.Id, Uid: usr.Id})
	if err != nil {
		ctx.JSON(http.StatusOK, Result{Code: 5, Msg: "系统错误"})
		a.l.Error("撤回为草稿失败", logger.Error(err))
		return
	}
	ctx.JSON(http.StatusOK, Result{Msg: "OK"})
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

// ListPubByAuthor 获取某作者的已发布文章列表
func (a *ArticleHandler) ListPubByAuthor(ctx *gin.Context, req AuthorPubListReq) (ginx.Result, error) {
	if req.Limit > 100 || req.Limit <= 0 {
		req.Limit = 10
	}
	idstr := ctx.Param("id")
	id, err := strconv.ParseInt(idstr, 10, 64)
	if err != nil {
		return ginx.Result{Code: 4, Msg: "参数错误"}, nil
	}

	authorOffset := req.Offset
	collected := make([]*articlev1.Article, 0, req.Limit)
	cursor := int32(0)
	const pageSize int32 = 100
	for int32(len(collected)) < req.Limit {
		pubResp, er := a.svc.ListPub(ctx, &articlev1.ListPubRequest{StartTime: timestamppb.Now(), Offset: cursor, Limit: pageSize})
		if er != nil {
			a.l.Error("获取推荐文章列表失败", logger.Error(er))
			return ginx.Result{Code: 5, Msg: "系统错误"}, er
		}
		arts := pubResp.GetArticles()
		if len(arts) == 0 {
			break
		}
		for _, art := range arts {
			if art.Author != nil && art.Author.Id == id && art.Status == 2 {
				if authorOffset > 0 {
					authorOffset--
					continue
				}
				collected = append(collected, art)
				if int32(len(collected)) >= req.Limit {
					break
				}
			}
		}
		if len(arts) < int(pageSize) || int32(len(collected)) >= req.Limit {
			break
		}
		cursor += pageSize
	}

	if len(collected) == 0 {
		return ginx.Result{Data: []ArticlePubVo{}}, nil
	}

	ids := make([]int64, 0, len(collected))
	for _, art := range collected {
		ids = append(ids, art.Id)
	}

	intrMap := make(map[int64]*intrv1.Interactive)
	intrResp, er := a.intrSvc.GetByIds(ctx, &intrv1.GetByIdsRequest{Biz: a.biz, Ids: ids})
	if er == nil && intrResp != nil {
		intrMap = intrResp.GetIntrs()
	}

	type userIntr struct{ liked, collected bool }
	userMap := make(map[int64]userIntr, len(ids))
	var uid int64
	if rawUser, ok := ctx.Get("user"); ok {
		if claims, ok2 := rawUser.(ginx.UserClaims); ok2 {
			uid = claims.Id
		}
	}
	if uid > 0 {
		for _, aid := range ids {
			r, e := a.intrSvc.Get(ctx, &intrv1.GetRequest{Biz: a.biz, BizId: aid, Uid: uid})
			if e != nil || r == nil || r.Intr == nil {
				continue
			}
			userMap[aid] = userIntr{liked: r.Intr.Liked, collected: r.Intr.Collected}
		}
	}

	result := make([]ArticlePubVo, 0, len(collected))
	for _, art := range collected {
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
		if intr, ok := intrMap[art.Id]; ok && intr != nil {
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

	return ginx.Result{Data: result}, nil
}

// ListCollected 获取用户收藏的文章列表
func (a *ArticleHandler) ListCollected(ctx *gin.Context, req CollectedListReq, uc ginx.UserClaims) (ginx.Result, error) {
	if req.Limit > 100 {
		req.Limit = 100
	}

	// 获取用户收藏的文章ID列表
	collectedResp, err := a.intrSvc.GetCollectedBizIds(ctx, &intrv1.GetCollectedBizIdsRequest{
		Biz:    a.biz,
		Uid:    uc.Id,
		Offset: req.Offset,
		Limit:  req.Limit,
	})
	if err != nil {
		a.l.Error("获取用户收藏列表失败", logger.Error(err))
		return ginx.Result{Code: 5, Msg: "系统错误"}, err
	}

	if len(collectedResp.BizIds) == 0 {
		return ginx.Result{Data: []ArticlePubVo{}, Msg: "OK"}, nil
	}

	// 获取文章详情（使用ListPub方法，按ID过滤）
	articles := make([]*articlev1.Article, 0, len(collectedResp.BizIds))
	for _, id := range collectedResp.BizIds {
		artResp, err := a.svc.GetPublishedById(ctx, &articlev1.GetPublishedByIdRequest{Id: id})
		if err != nil {
			a.l.Error("获取文章详情失败", logger.Error(err), logger.Int64("articleId", id))
			continue
		}
		if artResp.Article != nil && artResp.Article.Status == 2 {
			articles = append(articles, artResp.Article)
		}
	}

	// 获取互动数据
	intrMap := make(map[int64]*intrv1.Interactive)
	intrResp, err := a.intrSvc.GetByIds(ctx, &intrv1.GetByIdsRequest{Biz: a.biz, Ids: collectedResp.BizIds})
	if err == nil && intrResp != nil {
		intrMap = intrResp.GetIntrs()
	}

	// 构建返回结果
	result := make([]ArticlePubVo, 0, len(articles))
	for _, art := range articles {
		vo := ArticlePubVo{
			Id:         art.Id,
			Title:      art.Title,
			Abstract:   art.Abstract,
			CoverImage: art.CoverImage,
			Ctime:      art.Ctime.AsTime().Format(time.DateTime),
			Utime:      art.Utime.AsTime().Format(time.DateTime),
			Collected:  true, // 用户收藏的文章，肯定是已收藏状态
		}

		if art.Author != nil {
			vo.Author = AuthorVo{Id: art.Author.Id, Name: art.Author.Name, Avatar: art.Author.Avatar}
		}

		if intr, ok := intrMap[art.Id]; ok && intr != nil {
			vo.ReadCnt = intr.ReadCnt
			vo.LikeCnt = intr.LikeCnt
			vo.CollectCnt = intr.CollectCnt
			vo.Liked = intr.Liked
		}

		result = append(result, vo)
	}

	return ginx.Result{
		Data: result,
		Msg:  "OK",
	}, nil
}

func (a *ArticleHandler) ListFollowing(ctx *gin.Context, req FollowingListReq, uc ginx.UserClaims) (ginx.Result, error) {
	if req.Limit > 100 || req.Limit <= 0 {
		req.Limit = 10
	}
	frResp, err := a.followSvc.GetFollowee(ctx, &followv1.GetFolloweeRequest{Follower: uc.Id, Offset: 0, Limit: 1000})
	if err != nil {
		a.l.Error("获取关注列表失败", logger.Error(err))
		return ginx.Result{Code: 5, Msg: "系统错误"}, err
	}
	rels := frResp.GetFollowRelations()
	if len(rels) == 0 {
		return ginx.Result{Data: []ArticlePubVo{}}, nil
	}
	followeeSet := make(map[int64]struct{}, len(rels))
	for _, r := range rels {
		if r == nil {
			continue
		}
		followeeSet[r.GetFollowee()] = struct{}{}
	}

	articleOffset := req.Offset
	selected := make([]*articlev1.Article, 0, req.Limit)
	cursor := int32(0)
	const pageSize int32 = 100
	for int32(len(selected)) < req.Limit {
		pubResp, er := a.svc.ListPub(ctx, &articlev1.ListPubRequest{StartTime: timestamppb.Now(), Offset: cursor, Limit: pageSize})
		if er != nil {
			a.l.Error("获取推荐文章列表失败", logger.Error(er))
			return ginx.Result{Code: 5, Msg: "系统错误"}, er
		}
		arts := pubResp.GetArticles()
		if len(arts) == 0 {
			break
		}
		for _, art := range arts {
			if art.Author != nil {
				if _, ok := followeeSet[art.Author.Id]; ok && art.Status == 2 {
					if articleOffset > 0 {
						articleOffset--
						continue
					}
					selected = append(selected, art)
					if int32(len(selected)) >= req.Limit {
						break
					}
				}
			}
		}
		if len(arts) < int(pageSize) || int32(len(selected)) >= req.Limit {
			break
		}
		cursor += pageSize
	}

	if len(selected) == 0 {
		return ginx.Result{Data: []ArticlePubVo{}}, nil
	}

	ids := make([]int64, 0, len(selected))
	for _, art := range selected {
		ids = append(ids, art.Id)
	}

	intrMap := make(map[int64]*intrv1.Interactive)
	intrResp, er := a.intrSvc.GetByIds(ctx, &intrv1.GetByIdsRequest{Biz: a.biz, Ids: ids})
	if er == nil && intrResp != nil {
		intrMap = intrResp.GetIntrs()
	}

	type userIntr struct{ liked, collected bool }
	userMap := make(map[int64]userIntr, len(ids))
	var uid int64
	if rawUser, ok := ctx.Get("user"); ok {
		if claims, ok2 := rawUser.(ginx.UserClaims); ok2 {
			uid = claims.Id
		}
	}
	if uid > 0 {
		for _, aid := range ids {
			r, e := a.intrSvc.Get(ctx, &intrv1.GetRequest{Biz: a.biz, BizId: aid, Uid: uid})
			if e != nil || r == nil || r.Intr == nil {
				continue
			}
			userMap[aid] = userIntr{liked: r.Intr.Liked, collected: r.Intr.Collected}
		}
	}

	result := make([]ArticlePubVo, 0, len(selected))
	for _, art := range selected {
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
		if intr, ok := intrMap[art.Id]; ok && intr != nil {
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

	return ginx.Result{Data: result}, nil
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

type TagRankingReq struct {
	Tag    string `json:"tag" form:"tag"`
	Offset int32  `json:"offset" form:"offset"`
	Limit  int32  `json:"limit" form:"limit"`
}

type OfficialTagListReq struct {
	Tag    string `json:"tag" form:"tag"`
	Offset int32  `json:"offset" form:"offset"`
	Limit  int32  `json:"limit" form:"limit"`
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

func (a *ArticleHandler) ensureSearchClient() searchv1.SearchServiceClient {
	if a.searchSvc != nil {
		return a.searchSvc
	}
	type Config struct {
		Target string `json:"target"`
		Secure bool   `json:"secure"`
	}
	var cfg Config
	_ = viper.UnmarshalKey("grpc.client.search", &cfg)
	ecli, err := func() (*clientv3.Client, error) {
		var ecfg clientv3.Config
		if er := viper.UnmarshalKey("etcd", &ecfg); er != nil {
			return nil, er
		}
		return clientv3.New(ecfg)
	}()
	if err != nil {
		return nil
	}
	rs, err := resolver.NewBuilder(ecli)
	if err != nil {
		return nil
	}
	opts := []grpc.DialOption{grpc.WithResolvers(rs)}
	if !cfg.Secure {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}
	cc, err := grpc.Dial(cfg.Target, opts...)
	if err != nil {
		return nil
	}
	a.searchSvc = searchv1.NewSearchServiceClient(cc)
	return a.searchSvc
}

func (a *ArticleHandler) RankingByTag(ctx *gin.Context, req TagRankingReq) (ginx.Result, error) {
	if req.Tag == "" {
		return ginx.Result{Data: []ArticlePubVo{}}, nil
	}
	if req.Limit > 100 || req.Limit <= 0 {
		req.Limit = 10
	}
	sc := a.ensureSearchClient()
	if sc == nil {
		return ginx.Result{Code: 5, Msg: "系统错误"}, nil
	}
	sresp, err := sc.Search(ctx, &searchv1.SearchRequest{Expression: req.Tag, Uid: 0})
	if err != nil || sresp == nil || sresp.Article == nil {
		return ginx.Result{Data: []ArticlePubVo{}}, nil
	}
	arts := sresp.Article.Articles
	if len(arts) == 0 {
		return ginx.Result{Data: []ArticlePubVo{}}, nil
	}
	ids := make([]int64, 0, len(arts))
	for _, aitem := range arts {
		ids = append(ids, aitem.Id)
	}
	intrMap := make(map[int64]*intrv1.Interactive)
	intrResp, _ := a.intrSvc.GetByIds(ctx, &intrv1.GetByIdsRequest{Biz: a.biz, Ids: ids})
	if intrResp != nil {
		intrMap = intrResp.GetIntrs()
	}
	type item struct {
		id   int64
		like int64
		col  int64
	}
	arr := make([]item, 0, len(ids))
	for _, id := range ids {
		var lk, cl int64
		if intr := intrMap[id]; intr != nil {
			lk = intr.LikeCnt
			cl = intr.CollectCnt
		}
		arr = append(arr, item{id: id, like: lk, col: cl})
	}
	sort.Slice(arr, func(i, j int) bool {
		return arr[i].like+arr[i].col > arr[j].like+arr[j].col
	})
	start := int(req.Offset)
	end := start + int(req.Limit)
	if start >= len(arr) {
		return ginx.Result{Data: []ArticlePubVo{}}, nil
	}
	if end > len(arr) {
		end = len(arr)
	}
	page := arr[start:end]
	result := make([]ArticlePubVo, 0, len(page))
	var uid int64
	if raw, ok := ctx.Get("user"); ok {
		if claims, ok2 := raw.(ginx.UserClaims); ok2 {
			uid = claims.Id
		}
	}
	userMap := make(map[int64]struct{ liked, collected bool }, len(page))
	for _, it := range page {
		if uid > 0 {
			resp, er := a.intrSvc.Get(ctx, &intrv1.GetRequest{Biz: a.biz, BizId: it.id, Uid: uid})
			if er == nil && resp != nil && resp.Intr != nil {
				userMap[it.id] = struct{ liked, collected bool }{liked: resp.Intr.Liked, collected: resp.Intr.Collected}
			}
		}
		ar, er := a.svc.GetPublishedById(ctx, &articlev1.GetPublishedByIdRequest{Id: it.id})
		if er != nil || ar == nil || ar.Article == nil || ar.Article.Author == nil {
			continue
		}
		aid := ar.Article.Author.Id
		vo := ArticlePubVo{
			Id:         ar.Article.Id,
			Title:      ar.Article.Title,
			Abstract:   ar.Article.Abstract,
			CoverImage: ar.Article.CoverImage,
			Ctime:      formatProtoTime(ar.Article.Ctime),
			Utime:      formatProtoTime(ar.Article.Utime),
			ReadCnt: func() int64 {
				if intr := intrMap[ar.Article.Id]; intr != nil {
					return intr.ReadCnt
				}
				return 0
			}(),
			LikeCnt:    it.like,
			CollectCnt: it.col,
		}
		vo.Author = AuthorVo{Id: aid}
		if u, ok := userMap[ar.Article.Id]; ok {
			vo.Liked = u.liked
			vo.Collected = u.collected
		}
		result = append(result, vo)
	}
	return ginx.Result{Data: result}, nil
}

func (a *ArticleHandler) ListByOfficialTag(ctx *gin.Context, req OfficialTagListReq) (ginx.Result, error) {
	if req.Tag == "" {
		return ginx.Result{Data: []ArticlePubVo{}}, nil
	}
	if req.Limit > 100 || req.Limit <= 0 {
		req.Limit = 10
	}
	if a.tagSvc != nil {
		resp, _ := a.tagSvc.GetTags(ctx, &tagv1.GetTagsRequest{Uid: 0})
		if resp == nil || len(resp.Tag) == 0 {
			return ginx.Result{Data: []ArticlePubVo{}}, nil
		}
		found := false
		for _, t := range resp.Tag {
			if t != nil && t.Name == req.Tag {
				found = true
				break
			}
		}
		if !found {
			return ginx.Result{Data: []ArticlePubVo{}}, nil
		}
	}
	collected := make([]*articlev1.Article, 0, req.Limit)
	skip := req.Offset
	cursor := int32(0)
	const pageSize int32 = 100
	for int32(len(collected)) < req.Limit {
		pubResp, er := a.svc.ListPub(ctx, &articlev1.ListPubRequest{StartTime: timestamppb.Now(), Offset: cursor, Limit: pageSize})
		if er != nil {
			a.l.Error("获取推荐文章列表失败", logger.Error(er))
			return ginx.Result{Code: 5, Msg: "系统错误"}, er
		}
		arts := pubResp.GetArticles()
		if len(arts) == 0 {
			break
		}
		for _, art := range arts {
			if art.Status != 2 || art.Author == nil {
				continue
			}
			tgResp, _ := a.tagSvc.GetBizTags(ctx, &tagv1.GetBizTagsRequest{Uid: art.Author.Id, Biz: a.biz, BizId: art.Id})
			if tgResp == nil || len(tgResp.Tags) == 0 {
				continue
			}
			matched := false
			for _, t := range tgResp.Tags {
				if t != nil && t.Name == req.Tag {
					matched = true
					break
				}
			}
			if !matched {
				continue
			}
			if skip > 0 {
				skip--
				continue
			}
			collected = append(collected, art)
			if int32(len(collected)) >= req.Limit {
				break
			}
		}
		if len(arts) < int(pageSize) || int32(len(collected)) >= req.Limit {
			break
		}
		cursor += pageSize
	}
	if len(collected) == 0 {
		return ginx.Result{Data: []ArticlePubVo{}}, nil
	}
	ids := make([]int64, 0, len(collected))
	for _, aitem := range collected {
		ids = append(ids, aitem.Id)
	}
	intrMap := make(map[int64]*intrv1.Interactive)
	intrResp, _ := a.intrSvc.GetByIds(ctx, &intrv1.GetByIdsRequest{Biz: a.biz, Ids: ids})
	if intrResp != nil {
		intrMap = intrResp.GetIntrs()
	}
	result := make([]ArticlePubVo, 0, len(collected))
	var uid int64
	if raw, ok := ctx.Get("user"); ok {
		if claims, ok2 := raw.(ginx.UserClaims); ok2 {
			uid = claims.Id
		}
	}
	userMap := make(map[int64]struct{ liked, collected bool }, len(collected))
	for _, it := range collected {
		if uid > 0 {
			resp, er := a.intrSvc.Get(ctx, &intrv1.GetRequest{Biz: a.biz, BizId: it.Id, Uid: uid})
			if er == nil && resp != nil && resp.Intr != nil {
				userMap[it.Id] = struct{ liked, collected bool }{liked: resp.Intr.Liked, collected: resp.Intr.Collected}
			}
		}
		ar, er := a.svc.GetPublishedById(ctx, &articlev1.GetPublishedByIdRequest{Id: it.Id})
		if er != nil || ar == nil || ar.Article == nil {
			continue
		}
		vo := ArticlePubVo{
			Id:         ar.Article.Id,
			Title:      ar.Article.Title,
			Abstract:   ar.Article.Abstract,
			CoverImage: ar.Article.CoverImage,
			Ctime:      formatProtoTime(ar.Article.Ctime),
			Utime:      formatProtoTime(ar.Article.Utime),
		}
		if ar.Article.Author != nil {
			vo.Author = AuthorVo{Id: ar.Article.Author.Id, Name: ar.Article.Author.Name, Avatar: ar.Article.Author.Avatar}
		}
		if intr := intrMap[ar.Article.Id]; intr != nil {
			vo.ReadCnt = intr.ReadCnt
			vo.LikeCnt = intr.LikeCnt
			vo.CollectCnt = intr.CollectCnt
		}
		if u, ok := userMap[ar.Article.Id]; ok {
			vo.Liked = u.liked
			vo.Collected = u.collected
		}
		result = append(result, vo)
	}
	return ginx.Result{Data: result}, nil
}
func (a *ArticleHandler) OfficialTags(ctx *gin.Context) {
	if a.tagSvc == nil {
		ctx.JSON(http.StatusOK, Result{Data: []string{}})
		return
	}
	resp, err := a.tagSvc.GetTags(ctx, &tagv1.GetTagsRequest{Uid: 0})
	if err != nil {
		ctx.JSON(http.StatusOK, Result{Code: 5, Msg: "系统错误"})
		return
	}
	names := make([]string, 0, len(resp.GetTag()))
	for _, t := range resp.GetTag() {
		if t != nil && t.Name != "" {
			names = append(names, t.Name)
		}
	}
	ctx.JSON(http.StatusOK, Result{Data: names})
}
