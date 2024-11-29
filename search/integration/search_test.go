package integration

import (
	searchv1 "Webook/api/proto/gen/api/proto/search/v1"
	grpc2 "Webook/search/grpc"
	"Webook/search/integration/startup"
	"context"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type SearchTestSuite struct {
	suite.Suite
	searchSvc *grpc2.SearchServiceServer
	syncSvc   *grpc2.SyncServiceServer
}

func (s *SearchTestSuite) SetupSuite() {
	s.searchSvc = startup.InitSearchServer()
	s.syncSvc = startup.InitSyncServer()
}

func (s *SearchTestSuite) TestSearch() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	data, err := json.Marshal(BizTags{
		Uid:   1001,
		Biz:   "article",
		BizId: 123,
		Tags:  []string{"Jerry"},
	})
	require.NoError(s.T(), err)
	_, err = s.syncSvc.InputAny(ctx, &searchv1.InputAnyRequest{
		IndexName: "tags_index",
		DocId:     "abcd",
		Data:      string(data),
	})
	require.NoError(s.T(), err)

	_, err = s.syncSvc.InputUser(ctx, &searchv1.InputUserRequest{
		User: &searchv1.User{
			Id:       123,
			Nickname: "Tom White",
		},
	})
	require.NoError(s.T(), err)
	_, err = s.syncSvc.InputArticle(ctx, &searchv1.InputArticleRequest{
		Article: &searchv1.Article{
			Id:     123,
			Title:  "Tom 的小秘密",
			Status: 2,
		},
	})
	require.NoError(s.T(), err)
	_, err = s.syncSvc.InputArticle(ctx, &searchv1.InputArticleRequest{
		Article: &searchv1.Article{
			Id:      124,
			Content: "这是内容，Tom 的小秘密",
			Status:  2,
		},
	})
	require.NoError(s.T(), err)
	resp, err := s.searchSvc.Search(ctx, &searchv1.SearchRequest{
		Expression: "Tom 内容 Jerry",
		Uid:        1001,
	})
	require.NoError(s.T(), err)
	assert.Equal(s.T(), 1, len(resp.User.Users))
	assert.Equal(s.T(), 2, len(resp.Article.Articles))
}

type BizTags struct {
	Uid   int64    `json:"uid"`
	Biz   string   `json:"biz"`
	BizId int64    `json:"biz_id"`
	Tags  []string `json:"tags"`
}

func TestSearchService(t *testing.T) {
	suite.Run(t, new(SearchTestSuite))
}
