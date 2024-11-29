package dao

import (
	"context"
	"github.com/ecodeclub/ekit/slice"
	"github.com/goccy/go-json"
	"github.com/olivere/elastic/v7"
	"strconv"
	"strings"
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
	queryString := strings.Join(keywords, " ")
	//	需要把ids变成any的类型
	ids := slice.Map(tagArtIds, func(idx int, src int64) any {
		return src
	})
	//term 是字段查询，相当于去匹配字段的
	//match是全文查询，找对应的值的
	should := elastic.NewBoolQuery().Should(
		elastic.NewTermsQuery("id", ids...).Boost(2),
		elastic.NewMatchQuery("title", queryString),
		elastic.NewMatchQuery("content", queryString))
	query := elastic.NewBoolQuery().Must(
		should,
		elastic.NewTermQuery("status", "2"))
	resp, err := a.client.Search(ArticleIndexName).Query(query).Do(ctx)
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
