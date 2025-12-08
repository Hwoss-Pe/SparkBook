package web

import articlev1 "Webook/api/proto/gen/api/proto/article/v1"

type RewardReq struct {
	Id  int64 `json:"id"`
	Amt int64 `json:"amt"`
}

type LikeReq struct {
	Id   int64 `json:"id"`
	Like bool  `json:"like"`
}

type CollectReq struct {
	Id  int64 `json:"id"`
	Cid int64 `json:"cid"`
}

type ArticleVo struct {
	Id    int64  `json:"id"`
	Title string `json:"title"`
	// 摘要
	Abstract string `json:"abstract"`
	// 内容
	Content string `json:"content"`
	// 封面图片
	CoverImage string   `json:"coverImage"`
	Status     int32    `json:"status"`
	Author     AuthorVo `json:"author"`
	Ctime      string   `json:"ctime"`
	Utime      string   `json:"utime"`

	// 点赞之类的信息
	LikeCnt    int64 `json:"likeCnt"`
	CollectCnt int64 `json:"collectCnt"`
	ReadCnt    int64 `json:"readCnt"`

	// 个人是否点赞的信息
	Liked     bool `json:"liked"`
	Collected bool `json:"collected"`
}

type ArticleReq struct {
	Id         int64  `json:"id"`
	Title      string `json:"title"`
	Content    string `json:"content"`
	CoverImage string `json:"coverImage"`
}

func (req ArticleReq) toDTO(uid int64) *articlev1.Article {
	return &articlev1.Article{
		Id:         req.Id,
		Title:      req.Title,
		Content:    req.Content,
		CoverImage: req.CoverImage,
		Author: &articlev1.Author{
			Id: uid,
		},
	}
}

// AuthorVo 作者信息
type AuthorVo struct {
	Id     int64  `json:"id"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
}

// ArticlePubVo 推荐文章列表返回的文章信息
type ArticlePubVo struct {
	Id         int64    `json:"id"`
	Title      string   `json:"title"`
	Abstract   string   `json:"abstract"`
	CoverImage string   `json:"coverImage"`
	Author     AuthorVo `json:"author"`
	Ctime      string   `json:"ctime"`
	Utime      string   `json:"utime"`
	ReadCnt    int64    `json:"readCnt"`
	LikeCnt    int64    `json:"likeCnt"`
	CollectCnt int64    `json:"collectCnt"`
	Liked      bool     `json:"liked"`
	Collected  bool     `json:"collected"`
}
