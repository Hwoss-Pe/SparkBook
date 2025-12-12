package dao

import (
	"context"
	"strconv"
	"strings"

	"github.com/ecodeclub/ekit/slice"
	"github.com/goccy/go-json"
	"github.com/olivere/elastic/v7"
)

const ArticleIndexName = "article_index"
const TagIndexName = "tags_index"

type Article struct {
	Id      int64    `json:"id"`
	Title   string   `json:"title"`
	Status  int32    `json:"status"`
	Content string   `json:"content"`
	Tags    []string `json:"tags"`
}

type ArticleElasticDAO struct {
	client *elastic.Client
}

func (a *ArticleElasticDAO) InputArticle(ctx context.Context, article Article) error {
	_, err := a.client.Index().
		Index(ArticleIndexName).
		Id(strconv.FormatInt(article.Id, 10)).
		BodyJson(article).Do(ctx)
	return err
}

func (a *ArticleElasticDAO) Search(ctx context.Context, tagArtIds []int64, keywords []string) ([]Article, error) {
	// 当标签索引命中文章 ID 时，优先按 ID 精确过滤 + 状态
	if len(tagArtIds) > 0 {
		ids := slice.Map(tagArtIds, func(idx int, src int64) any { return src })
		query := elastic.NewBoolQuery().Must(
			elastic.NewTermsQuery("id", ids...),
			elastic.NewTermQuery("status", 2),
		)
		resp, err := a.client.Search(ArticleIndexName).Query(query).Size(100).Do(ctx)
		if err != nil {
			return nil, err
		}
		articles := make([]Article, 0, len(resp.Hits.Hits))
		for _, hit := range resp.Hits.Hits {
			var ele Article
			err = json.Unmarshal(hit.Source, &ele)
			articles = append(articles, ele)
		}
		return articles, nil
	}

	// 否则走全文检索（标题/内容），并要求至少命中一项
	queryString := strings.Join(keywords, " ")
	should := elastic.NewBoolQuery().Should(
		elastic.NewMatchQuery("title", queryString),
		elastic.NewMatchQuery("content", queryString),
	).MinimumShouldMatch("1")
	query := elastic.NewBoolQuery().Must(
		should,
		elastic.NewTermQuery("status", 2),
	)
	resp, err := a.client.Search(ArticleIndexName).Query(query).Size(100).Do(ctx)
	if err != nil {
		return nil, err
	}
	articles := make([]Article, 0, len(resp.Hits.Hits))
	for _, hit := range resp.Hits.Hits {
		var ele Article
		err = json.Unmarshal(hit.Source, &ele)
		articles = append(articles, ele)
	}
	return articles, nil
}

func NewArticleElasticDAO(client *elastic.Client) ArticleDAO {
	return &ArticleElasticDAO{client: client}
}
