package repository

import (
	"Webook/article/domain"
	"Webook/article/repository/cache"
	"Webook/article/repository/dao"
	"Webook/pkg/logger"
	"context"
	"github.com/ecodeclub/ekit/slice"
	"gorm.io/gorm"
	"time"
)

//go:generate mockgen -source=./type.go -package=repomocks -destination=mocks/article.mock.go ArticleRepository
type ArticleRepository interface {
	Create(ctx context.Context, art domain.Article) (int64, error)
	Update(ctx context.Context, art domain.Article) error
	List(ctx context.Context, author int64,
		offset int, limit int) ([]domain.Article, error)

	// Sync 本身要求先保存到制作库，再同步到线上库
	Sync(ctx context.Context, art domain.Article) (int64, error)
	// SyncStatus 仅仅同步状态
	SyncStatus(ctx context.Context, uid, id int64, status domain.ArticleStatus) error
	GetById(ctx context.Context, id int64) (domain.Article, error)

	GetPublishedById(ctx context.Context, id int64) (domain.Article, error)
	ListPub(ctx context.Context, utime time.Time, offset int, limit int) ([]domain.Article, error)
}

type CachedArticleRepository struct {
	// 操作单一的库
	dao   dao.ArticleDAO
	cache cache.ArticleCache

	//// SyncV1 用
	//authorDAO dao.ArticleAuthorDAO
	//readerDAO dao.ArticleReaderDAO

	// SyncV2 用
	db *gorm.DB
	l  logger.Logger
}

func NewArticleRepository(dao dao.ArticleDAO, cache cache.ArticleCache, db *gorm.DB, l logger.Logger) ArticleRepository {
	return &CachedArticleRepository{dao: dao, cache: cache, db: db, l: l}
}

func (c *CachedArticleRepository) Create(ctx context.Context, art domain.Article) (int64, error) {
	id, err := c.dao.Insert(ctx, c.toEntity(art))
	if err != nil {
		return 0, err
	}
	author := art.Author.Id
	//在用户发布新的文章的时候需要把第一页缓存记得删除
	err = c.cache.DelFirstPage(ctx, author)
	if err != nil {
		c.l.Error("删除缓存失败",
			logger.Int64("author", author), logger.Error(err))
	}
	return id, nil
}

func (c *CachedArticleRepository) Update(ctx context.Context, art domain.Article) error {
	err := c.dao.UpdateById(ctx, c.toEntity(art))
	if err != nil {
		return err
	}
	author := art.Author.Id
	err = c.cache.DelFirstPage(ctx, author)
	if err != nil {
		c.l.Error("删除缓存失败",
			logger.Int64("author", author), logger.Error(err))
	}
	return nil
}

func (c *CachedArticleRepository) List(ctx context.Context, author int64, offset int, limit int) ([]domain.Article, error) {
	// 只有第一页才走缓存，并且假定一页只有 100 条
	// 也就是说，如果前端允许创作者调整页的大小
	// 那么只有 100 这个页大小这个默认情况下，会走索引
	if offset == 0 && limit <= 100 {
		page, err := c.cache.GetFirstPage(ctx, author)
		//这是走缓存的快路径
		if err == nil {
			//如果他是这样我就提前准备缓存
			if err == nil {
				go func() {
					c.preCache(ctx, page)
				}()
				return page, nil
			}
		}
	}
	// 慢路径
	arts, err := c.dao.GetByAuthor(ctx, author, offset, limit)
	if err != nil {
		return nil, err
	}
	res := slice.Map[dao.Article, domain.Article](arts,
		func(idx int, src dao.Article) domain.Article {
			return c.ToDomain(src)
		})
	go func() {
		c.preCache(ctx, res)
	}()
	//	重新把第一页缓存设置回去
	err = c.cache.SetFirstPage(ctx, author, res)
	if err != nil {
		c.l.Error("刷新第一页文章的缓存失败",
			logger.Int64("author", author), logger.Error(err))
	}
	return res, nil
}

func (c *CachedArticleRepository) Sync(ctx context.Context, art domain.Article) (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (c *CachedArticleRepository) SyncStatus(ctx context.Context,
	uid, id int64, status domain.ArticleStatus) error {
	return c.dao.SyncStatus(ctx, uid, id, status.ToUint8())
}

func (c *CachedArticleRepository) GetById(ctx context.Context, id int64) (domain.Article, error) {
	cachedArt, err := c.cache.Get(ctx, id)
	if err == nil {
		return cachedArt, nil
	}
	art, err := c.dao.GetById(ctx, id)
	if err != nil {
		return domain.Article{}, err
	}
	//在这里并没有回写到redis
	return c.ToDomain(art), nil
}

func (c *CachedArticleRepository) GetPublishedById(ctx context.Context, id int64) (domain.Article, error) {
	res, err := c.cache.GetPub(ctx, id)
	if err == nil {
		return res, err
	}
	art, err := c.dao.GetPubById(ctx, id)
	if err != nil {
		return domain.Article{}, err
	}
	res = domain.Article{
		Id:      art.Id,
		Title:   art.Title,
		Status:  domain.ArticleStatus(art.Status),
		Content: art.Content,
	}
	//	异步设置缓存 也可以同步
	go func() {
		if err = c.cache.SetPub(ctx, res); err != nil {
			c.l.Error("缓存已发表文章失败",
				logger.Error(err), logger.Int64("aid", res.Id))
		}
	}()
	return res, nil
}

func (c *CachedArticleRepository) ListPub(ctx context.Context, utime time.Time, offset int, limit int) ([]domain.Article, error) {
	articles, err := c.dao.ListPubByUtime(ctx, utime, offset, limit)
	if err != nil {
		return nil, err
	}
	return slice.Map[dao.PublishedArticle, domain.Article](articles, func(idx int, src dao.PublishedArticle) domain.Article {
		return c.ToDomain(dao.Article(src))
	}), nil
}

func (c *CachedArticleRepository) ToDomain(art dao.Article) domain.Article {
	return domain.Article{
		Id:      art.Id,
		Title:   art.Title,
		Status:  domain.ArticleStatus(art.Status),
		Content: art.Content,
		Author: domain.Author{
			Id: art.AuthorId,
		},
	}
}

func (c *CachedArticleRepository) toEntity(art domain.Article) dao.Article {
	return dao.Article{
		Id:       art.Id,
		Title:    art.Title,
		Content:  art.Content,
		AuthorId: art.Author.Id,
		// 这一步，就是将领域状态转化为存储状态。
		Status: uint8(art.Status),
	}
}

// 我默认用户进作者简介的时候只会点击第一个文章
func (c *CachedArticleRepository) preCache(ctx context.Context, arts []domain.Article) {
	// 1MB
	const contentSizeThreshold = 1024 * 1024

	if len(arts) > 0 && len(arts[0].Content) <= contentSizeThreshold {
		// 准备缓存1Mb的
		if err := c.cache.Set(ctx, arts[0]); err != nil {
			c.l.Error("提前准备缓存失败", logger.Error(err))
		}
	}
}
