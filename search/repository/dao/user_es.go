package dao

import (
	"context"
	"encoding/json"
	"github.com/olivere/elastic/v7"
	"strconv"
	"strings"
)

const UserIndexName = "user_index"

type User struct {
	Id       int64  `json:"id"`
	Email    string `json:"email"`
	Nickname string `json:"nickname"`
	Phone    string `json:"phone"`
	Avatar   string `json:"avatar"`
}

type UserElasticDAO struct {
	client *elastic.Client
}

func (u *UserElasticDAO) InputUser(ctx context.Context, user User) error {
	_, err := u.client.Index().
		Index(UserIndexName).
		Id(strconv.FormatInt(user.Id, 10)).
		BodyJson(user).Do(ctx)
	return err
}

func (u *UserElasticDAO) Search(ctx context.Context, keywords []string) ([]User, error) {
	// 假定上面传入的 keywords 是经过了处理的
	queryString := strings.Join(keywords, " ")
	query := elastic.NewBoolQuery().Must(elastic.NewMatchQuery("nickname", queryString))
	resp, err := u.client.Search().Index(UserIndexName).Query(query).Do(ctx)
	if err != nil {
		return nil, err
	}
	users := make([]User, 0, len(resp.Hits.Hits))
	for _, hit := range resp.Hits.Hits {
		var user User
		err := json.Unmarshal(hit.Source, &user)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func NewUserElasticDAO(client *elastic.Client) UserDAO {
	return &UserElasticDAO{
		client: client,
	}
}
