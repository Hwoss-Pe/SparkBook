package service

import (
	"Webook/article/domain"
	"Webook/article/events"
	"Webook/article/repository"
	"Webook/article/repository/dao"
	"Webook/pkg/logger"
	"context"
	"time"

	"golang.org/x/sync/errgroup"
)

type ArticleService interface {
	Save(ctx context.Context, art domain.Article) (int64, error)
	Publish(ctx context.Context, art domain.Article) (int64, error)
	Withdraw(ctx context.Context, uid, id int64) error
	Unpublish(ctx context.Context, uid, id int64) error
	PublishV1(ctx context.Context, art domain.Article) (int64, error)
	List(ctx context.Context, author int64, offset, limit int) ([]domain.Article, error)
	GetById(ctx context.Context, id int64) (domain.Article, error)
	// GetPublishedById 查找已经发表的
	GetPublishedById(ctx context.Context, id, uid int64) (domain.Article, error)
	// ListPub 根据更新时间来分页，更新时间必须小于 startTime
	ListPub(ctx context.Context, startTime time.Time, offset, limit int) ([]domain.Article, error)
}

type articleService struct {
	// 1. 在 service 这一层使用两个 repository
	authorRepo repository.ArticleAuthorRepository
	readerRepo repository.ArticleReaderRepository
	userRepo   repository.AuthorRepository
	logger     logger.Logger

	repo repository.ArticleRepository

	// 搞个异步的
	producer events.Producer
}

func NewArticleService(logger logger.Logger, repo repository.ArticleRepository,
	authorRepo repository.ArticleAuthorRepository, readerRepo repository.ArticleReaderRepository,
	userRepo repository.AuthorRepository, producer events.Producer) ArticleService {
	return &articleService{
		logger:     logger,
		repo:       repo,
		authorRepo: authorRepo,
		readerRepo: readerRepo,
		userRepo:   userRepo,
		producer:   producer,
	}
}

func (a *articleService) Save(ctx context.Context, art domain.Article) (int64, error) {
	// 设置为未发表
	art.Status = domain.ArticleStatusUnpublished
	if art.Id > 0 {
		err := a.update(ctx, art)
		return art.Id, err
	}
	return a.create(ctx, art)
}

func (a *articleService) Publish(ctx context.Context, art domain.Article) (int64, error) {
	art.Status = domain.ArticleStatusPublished
	return a.repo.Sync(ctx, art)
}

func (a *articleService) Withdraw(ctx context.Context, uid, id int64) error {
	return a.repo.SyncStatus(ctx, uid, id, domain.ArticleStatusPrivate)
}

func (a *articleService) Unpublish(ctx context.Context, uid, id int64) error {
	art, err := a.repo.GetById(ctx, id)
	if err != nil {
		return err
	}
	if art.Author.Id != uid {
		return dao.ErrPossibleIncorrectAuthor
	}
	if art.Status == domain.ArticleStatusUnpublished {
		err = a.repo.DeleteDraft(ctx, uid, id)
		if err == nil {
			go func() {
				_, _ = a.repo.List(ctx, uid, 0, 100)
			}()
		}
		return err
	}
	return a.repo.SyncStatus(ctx, uid, id, domain.ArticleStatusUnpublished)
}

func (a *articleService) PublishV1(ctx context.Context, art domain.Article) (int64, error) {
	var (
		id  = art.Id
		err error
	)
	if art.Id == 0 {
		id, err = a.authorRepo.Create(ctx, art)
	} else {
		err = a.authorRepo.Update(ctx, art)
	}
	if err != nil {
		return 0, err
	}
	//	保证线上库和制作库的id是一样的
	art.Id = id
	for i := 0; i < 3; i++ {
		//可以重试三次
		err := a.readerRepo.Save(ctx, art)
		if err == nil {
			break
		}
		a.logger.Error("部分失败：保存数据到线上库失败",
			logger.Field{Key: "art_id", Value: id},
			logger.Error(err))
	}
	if err != nil {
		a.logger.Error("部分失败：保存数据到线上库重试都失败了",
			logger.Field{Key: "art_id", Value: id},
			logger.Error(err))
		return 0, err
	}
	return id, nil
}

func (a *articleService) List(ctx context.Context, author int64,
	offset, limit int) ([]domain.Article, error) {
	return a.repo.List(ctx, author, offset, limit)
}

func (a *articleService) GetById(ctx context.Context, id int64) (domain.Article, error) {
	return a.repo.GetById(ctx, id)
}

func (a *articleService) GetPublishedById(ctx context.Context, id, uid int64) (domain.Article, error) {
	var eg errgroup.Group
	var art *domain.Article
	var author *domain.Author
	var err error
	eg.Go(func() error {
		res, eerr := a.repo.GetPublishedById(ctx, id)
		art = &res
		return eerr
	})
	eg.Go(func() error {
		res, eerr := a.userRepo.FindAuthor(ctx, id)
		author = &res
		return eerr
	})
	if err = eg.Wait(); err != nil {
		return domain.Article{}, err
	}
	art.Author = *author
	res := *art
	//	然后还要在这里提高阅读量发送kafka
	go func() {
		if err == nil {
			er := a.producer.ProduceReadEvent(events.ReadEvent{
				Aid: id,
				Uid: uid,
			})
			if er != nil {
				a.logger.Error("发送消息失败",
					logger.Int64("uid", uid),
					logger.Int64("aid", id),
					logger.Error(err))
			}
		}
	}()
	return res, err
}

func (a *articleService) ListPub(ctx context.Context, startTime time.Time, offset, limit int) ([]domain.Article, error) {
	articles, err := a.repo.ListPub(ctx, startTime, offset, limit)
	if err != nil {
		return nil, err
	}

	// 批量获取作者信息
	for i := range articles {
		author, err := a.userRepo.FindAuthor(ctx, articles[i].Id)
		if err != nil {
			// 如果获取作者信息失败，记录日志但不影响整体流程
			a.logger.Error("获取作者信息失败",
				logger.Int64("articleId", articles[i].Id),
				logger.Int64("authorId", articles[i].Author.Id),
				logger.Error(err))
			continue
		}
		articles[i].Author = author
	}

	return articles, nil
}

func (a *articleService) create(ctx context.Context,
	art domain.Article) (int64, error) {
	return a.repo.Create(ctx, art)
}
func (a *articleService) update(ctx context.Context,
	art domain.Article) error {
	return a.repo.Update(ctx, art)
}
