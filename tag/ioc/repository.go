package ioc

import (
	"Webook/pkg/logger"
	"Webook/tag/repository"
	"Webook/tag/repository/cache"
	"Webook/tag/repository/dao"
	"context"
	"time"
)

func InitRepository(d dao.TagDAO, c cache.TagCache, l logger.Logger) repository.TagRepository {
	repo := repository.NewTagRepository(d, c, l)
	go func() {
		_, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()
		// 也可以同步执行。但是在一些场景下，同步执行会占用很长的时间，所以可以考虑异步执行。
		//repo.PreloadUserTags(ctx)
	}()
	return repo
}
