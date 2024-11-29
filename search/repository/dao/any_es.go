package dao

import "github.com/olivere/elastic/v7"

type AnyESDAO struct {
	client *elastic.Client
}
