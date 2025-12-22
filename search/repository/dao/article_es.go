package dao

import (
	"bytes"
	"context"
	"math"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/ecodeclub/ekit/slice"
	"github.com/goccy/go-json"
	"github.com/olivere/elastic/v7"
	"golang.org/x/sync/errgroup"
)

const ArticleIndexName = "article_index"
const TagIndexName = "tags_index"

type Article struct {
	Id         int64     `json:"id"`
	Title      string    `json:"title"`
	Status     int32     `json:"status"`
	Content    string    `json:"content"`
	Tags       []string  `json:"tags"`
	TitleVec   []float32 `json:"title_vec"`
	ContentVec []float32 `json:"content_vec"`
}

type ArticleElasticDAO struct {
	client *elastic.Client
}

func (a *ArticleElasticDAO) InputArticle(ctx context.Context, article Article) error {
	if len(article.TitleVec) == 0 {
		article.TitleVec = remoteEmb(article.Title)
	}
	if len(article.ContentVec) == 0 {
		article.ContentVec = remoteEmb(article.Content)
	}
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

	qs := strings.Join(keywords, " ")

	// 1. 准备并发查询
	var (
		eg       errgroup.Group
		vecArts  []Article
		textArts []Article
	)

	// 2. 向量检索 (Semantic Search)
	eg.Go(func() error {
		qv := remoteEmb(qs)
		params := map[string]interface{}{
			"query_vector": qv,
		}
		src := "double s1 = cosineSimilarity(params.query_vector, doc['title_vec']); double s2 = cosineSimilarity(params.query_vector, doc['content_vec']); return Math.max(s1, s2) + 1.0;"
		script := elastic.NewScript(src).Params(params)
		filter := elastic.NewBoolQuery().Must(elastic.NewTermQuery("status", 2))
		ssq := elastic.NewScriptScoreQuery(filter, script)
		resp, err := a.client.Search(ArticleIndexName).Query(ssq).Size(50).Do(ctx)
		if err != nil {
			return err
		}
		vecArts = make([]Article, 0, len(resp.Hits.Hits))
		for _, hit := range resp.Hits.Hits {
			var ele Article
			err = json.Unmarshal(hit.Source, &ele)
			if err == nil {
				vecArts = append(vecArts, ele)
			}
		}
		return nil
	})

	// 3. 文本检索 (BM25 Search)
	eg.Go(func() error {
		// 使用 MultiMatch 在标题和内容中查找
		query := elastic.NewBoolQuery().Must(
			elastic.NewMultiMatchQuery(qs, "title", "content"),
			elastic.NewTermQuery("status", 2),
		)
		resp, err := a.client.Search(ArticleIndexName).Query(query).Size(50).Do(ctx)
		if err != nil {
			return err
		}
		textArts = make([]Article, 0, len(resp.Hits.Hits))
		for _, hit := range resp.Hits.Hits {
			var ele Article
			err = json.Unmarshal(hit.Source, &ele)
			if err == nil {
				textArts = append(textArts, ele)
			}
		}
		return nil
	})

	// 4. 等待结果并处理 RRF
	if err := eg.Wait(); err != nil {
		return nil, err
	}

	// 5. RRF 融合 (Reciprocal Rank Fusion)
	// score = 1.0 / (k + rank)
	const k = 60.0
	type rrfScore struct {
		art   Article
		score float64
	}
	merged := make(map[int64]*rrfScore)

	// 处理向量结果
	for i, art := range vecArts {
		rank := float64(i + 1)
		score := 1.0 / (k + rank)
		if val, ok := merged[art.Id]; ok {
			val.score += score
		} else {
			merged[art.Id] = &rrfScore{art: art, score: score}
		}
	}

	// 处理文本结果
	for i, art := range textArts {
		rank := float64(i + 1)
		score := 1.0 / (k + rank)
		if val, ok := merged[art.Id]; ok {
			val.score += score
		} else {
			merged[art.Id] = &rrfScore{art: art, score: score}
		}
	}

	// 6. 排序
	res := make([]*rrfScore, 0, len(merged))
	for _, v := range merged {
		res = append(res, v)
	}
	// 分数降序
	sort.Slice(res, func(i, j int) bool {
		return res[i].score > res[j].score
	})

	// 7. 返回结果
	articles := make([]Article, 0, len(res))
	for _, r := range res {
		articles = append(articles, r.art)
	}
	return articles, nil
}

func NewArticleElasticDAO(client *elastic.Client) ArticleDAO {
	return &ArticleElasticDAO{client: client}
}

var embedURL = func() string {
	if u := os.Getenv("EMBED_URL"); u != "" {
		return u
	}
	return "http://localhost:8088/embed"
}()

func emb(s string) []float32 {
	const dims = 512
	vec := make([]float32, dims)
	if s == "" {
		return vec
	}
	tokens := strings.Fields(strings.ToLower(s))
	if len(tokens) == 0 {
		return vec
	}
	for _, t := range tokens {
		h := fnv32(t)
		idx := int(h % dims)
		vec[idx] += 1
	}
	var sum float64
	for i := 0; i < dims; i++ {
		sum += float64(vec[i] * vec[i])
	}
	norm := float32(math.Sqrt(sum))
	if norm == 0 {
		return vec
	}
	for i := 0; i < dims; i++ {
		vec[i] = vec[i] / norm
	}
	return vec
}

func remoteEmb(s string) []float32 {
	type req struct {
		Texts     []string `json:"texts"`
		Normalize bool     `json:"normalize"`
		IsQuery   bool     `json:"is_query"`
	}
	type resp struct {
		Vectors [][]float64 `json:"vectors"`
	}
	r := req{
		Texts:     []string{s},
		Normalize: true,
		IsQuery:   true,
	}
	bs, err := json.Marshal(r)
	if err != nil {
		return emb(s)
	}
	cli := &http.Client{Timeout: 3 * time.Second}
	respRaw, err := cli.Post(embedURL, "application/json", bytes.NewReader(bs))
	if err != nil {
		return emb(s)
	}
	defer respRaw.Body.Close()
	var rp resp
	dec := json.NewDecoder(respRaw.Body)
	if err = dec.Decode(&rp); err != nil {
		return emb(s)
	}
	if len(rp.Vectors) == 0 {
		return emb(s)
	}
	fv := make([]float32, len(rp.Vectors[0]))
	for i, v := range rp.Vectors[0] {
		fv[i] = float32(v)
	}
	return fv
}
func fnv32(s string) uint32 {
	const (
		offset32 = 2166136261
		prime32  = 16777619
	)
	var h uint32 = offset32
	for i := 0; i < len(s); i++ {
		h ^= uint32(s[i])
		h *= prime32
	}
	return h
}
