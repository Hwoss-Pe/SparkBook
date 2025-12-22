package ioc

import (
	articlev1 "Webook/api/proto/gen/api/proto/article/v1"
	followv1 "Webook/api/proto/gen/api/proto/follow/v1"
	intrv1 "Webook/api/proto/gen/api/proto/intr/v1"
	rankingv1 "Webook/api/proto/gen/api/proto/ranking/v1"
	rewardv1 "Webook/api/proto/gen/api/proto/reward/v1"
	searchv1 "Webook/api/proto/gen/api/proto/search/v1"
	tagv1 "Webook/api/proto/gen/api/proto/tag/v1"
	"Webook/bff/web"
	"Webook/pkg/logger"

	"github.com/spf13/viper"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/resolver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func InitArticleClient(ecli *clientv3.Client) articlev1.ArticleServiceClient {
	type Config struct {
		Target string `json:"target"`
		Secure bool   `json:"secure"`
	}
	var cfg Config
	err := viper.UnmarshalKey("grpc.client.article", &cfg)
	if err != nil {
		panic(err)
	}
	rs, err := resolver.NewBuilder(ecli)
	if err != nil {
		panic(err)
	}
	opts := []grpc.DialOption{grpc.WithResolvers(rs)}
	if !cfg.Secure {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}
	cc, err := grpc.Dial(cfg.Target, opts...)
	if err != nil {
		panic(err)
	}
	return articlev1.NewArticleServiceClient(cc)
}

func NewArticleHandler(artSvc articlev1.ArticleServiceClient,
	intrSvc intrv1.InteractiveServiceClient,
	rankingSvc rankingv1.RankingServiceClient,
	rewardSvc rewardv1.RewardServiceClient,
	l logger.Logger,
	followSvc followv1.FollowServiceClient,
	tagSvc tagv1.TagServiceClient,
	searchSvc searchv1.SearchServiceClient) *web.ArticleHandler {
	return web.NewArticleHandler(artSvc, intrSvc, rankingSvc, rewardSvc, l, followSvc, tagSvc, searchSvc)
}

func InitTagClient(ecli *clientv3.Client) tagv1.TagServiceClient {
	type Config struct {
		Target string `json:"target"`
		Secure bool   `json:"secure"`
	}
	var cfg Config
	err := viper.UnmarshalKey("grpc.client.tag", &cfg)
	if err != nil {
		panic(err)
	}
	rs, err := resolver.NewBuilder(ecli)
	if err != nil {
		panic(err)
	}
	opts := []grpc.DialOption{grpc.WithResolvers(rs)}
	if !cfg.Secure {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}
	cc, err := grpc.Dial(cfg.Target, opts...)
	if err != nil {
		panic(err)
	}
	return tagv1.NewTagServiceClient(cc)
}
