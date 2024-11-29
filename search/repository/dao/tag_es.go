package dao

import "github.com/olivere/elastic/v7"

type TagESDAO struct {
	client *elastic.Client
}
