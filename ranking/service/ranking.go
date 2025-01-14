package service

import (
	articlev1 "Webook/api/proto/gen/api/proto/article/v1"
	intrv1 "Webook/api/proto/gen/api/proto/intr/v1"
	"Webook/ranking/domain"
	"Webook/ranking/repository"
	"errors"
	"github.com/ecodeclub/ekit/queue"
	"github.com/ecodeclub/ekit/slice"
	"golang.org/x/net/context"
	"google.golang.org/protobuf/types/known/timestamppb"
	"math"
	"time"
)

//go:generate mockgen -source=./ranking.go -package=svcmocks -destination=./mocks/ranking.mock.go RankingService
type RankingService interface {
	// RankTopN 计算 TopN
	RankTopN(ctx context.Context) error
	// TopN 返回业务的 ID
	TopN(ctx context.Context) ([]domain.Article, error)
}

// BatchRankingService 分批计算
type BatchRankingService struct {
	intrSvc intrv1.InteractiveServiceClient
	artSvc  articlev1.ArticleServiceClient
	repo    repository.RankingRepository

	BatchSize int
	N         int

	scoreFunc func(likeCnt int64, utime time.Time) float64
}

func NewBatchRankingService(intrSvc intrv1.InteractiveServiceClient, artSvc articlev1.ArticleServiceClient,
	repo repository.RankingRepository) RankingService {
	res := &BatchRankingService{
		intrSvc:   intrSvc,
		artSvc:    artSvc,
		repo:      repo,
		N:         100,
		BatchSize: 100,
	}
	res.scoreFunc = res.score
	return res
}

func (a *BatchRankingService) RankTopN(ctx context.Context) error {
	arts, err := a.rankTopN(ctx)
	if err != nil {
		return err
	}
	// 准备放到缓存里面
	return a.repo.ReplaceTopN(ctx, arts)
}

func (a *BatchRankingService) TopN(ctx context.Context) ([]domain.Article, error) {
	return a.repo.GetTopN(ctx)
}

// 分数计算算法，依据的是它的点赞量和时间衰减进行计算
func (a *BatchRankingService) score(likeCnt int64, utime time.Time) float64 {
	const factor = 1.5
	return float64(likeCnt-1) / math.Pow(time.Since(utime).Hours()+2, factor)
}
func articleToDomain(article *articlev1.Article) domain.Article {
	domainArticle := domain.Article{}
	if article != nil {
		domainArticle.Id = article.GetId()
		domainArticle.Title = article.GetTitle()
		domainArticle.Status = domain.ArticleStatus(article.Status)
		domainArticle.Content = article.Content
		domainArticle.Author = domain.Author{
			Id:   article.GetAuthor().GetId(),
			Name: article.GetAuthor().GetName(),
		}
		domainArticle.Ctime = article.Ctime.AsTime()
		domainArticle.Utime = article.Utime.AsTime()
	}
	return domainArticle
}
func (a *BatchRankingService) rankTopN(ctx context.Context) ([]domain.Article, error) {
	//只计算七天内的，超过七天就不算热榜，这里设置一个utime进行计算
	now := time.Now()
	ddl := now.Add(-time.Hour * 24 * 7)
	offset := 0
	type Score struct {
		art   domain.Article
		score float64
	}
	//	用一个优先队列维持住topN的id
	priorityQueue := queue.NewPriorityQueue[Score](a.N, func(src Score, dst Score) int {
		//原来的比目的小就返回-1
		if src.score > dst.score {
			return 1
		} else if src.score == dst.score {
			return 0
		} else {
			return -1
		}
	})
	for {
		arts, err := a.artSvc.ListPub(ctx, &articlev1.ListPubRequest{
			StartTime: timestamppb.New(now),
			Offset:    int32(offset),
			Limit:     int32(a.BatchSize),
		})
		if err != nil {
			return nil, err
		}
		//	转化成domain
		dms := make([]domain.Article, 0, len(arts.Articles))
		for _, art := range arts.Articles {
			dms = append(dms, articleToDomain(art))
		}
		ids := slice.Map[domain.Article, int64](dms, func(idx int, src domain.Article) int64 {
			return src.Id
		})
		idsResponses, err := a.intrSvc.GetByIds(ctx, &intrv1.GetByIdsRequest{
			Biz: "article", Ids: ids,
		})
		if err != nil {
			return nil, err
		}
		minScore := float64(0)
		for _, art := range dms {
			interactive, ok := idsResponses.GetIntrs()[art.Id]
			if !ok {
				continue
			}
			score := a.scoreFunc(interactive.LikeCnt, art.Utime)
			if score > minScore {
				ele := Score{
					art:   art,
					score: score,
				}
				err := priorityQueue.Enqueue(ele)
				//	如果队列满了
				if errors.Is(err, queue.ErrOutOfCapacity) {
					_, _ = priorityQueue.Dequeue()
					err = priorityQueue.Enqueue(ele)
				}
			} else {
				minScore = score
			}
		}
		//不够一批或者时间太久远的就遍历结束
		if len(dms) == 0 || len(dms) < a.BatchSize || dms[len(dms)-1].Utime.Before(ddl) {
			break
		}
		offset = offset + len(dms)
	}
	ql := priorityQueue.Len()
	res := make([]domain.Article, ql)
	//队列里面的是从高到低的,因此弄成低到高的
	for i := ql - 1; i >= 0; i-- {
		val, _ := priorityQueue.Dequeue()
		res[i] = val.art
	}
	return res, nil
}
