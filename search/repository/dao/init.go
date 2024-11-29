package dao

import (
	"context"
	_ "embed"
	"fmt"
	"github.com/olivere/elastic/v7"
	"golang.org/x/sync/errgroup"
	"time"
)

var (
	//go:embed indexs/user_index.json
	userIndex string
	//go:embed indexs/article_index.json
	articleIndex string
	//go:embed indexs/tags_index.json
	tagIndex string
)

func InitES(client *elastic.Client) error {
	const timeout = time.Second * 10
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	var eg errgroup.Group
	eg.Go(func() error {
		return tryCreateIndex(ctx, client, UserIndexName, userIndex)
	})
	eg.Go(func() error {
		return tryCreateIndex(ctx, client, ArticleIndexName, articleIndex)
	})
	eg.Go(func() error { return tryCreateIndex(ctx, client, TagIndexName, tagIndex) })

	return eg.Wait()
}

func tryCreateIndex(ctx context.Context, client *elastic.Client, idxName, idxCfg string) error {
	ok, err := client.IndexExists(idxName).Do(ctx)
	if err != nil {
		return fmt.Errorf("检测 %s 是否存在失败 %w", idxName, err)
	}
	if ok {
		return nil
	}
	_, err = client.CreateIndex(idxName).Body(idxCfg).Do(ctx)

	if err != nil {
		return fmt.Errorf("初始化 %s 失败 %w", idxName, err)
	}
	return err
}
