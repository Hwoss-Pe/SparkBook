package dao

import (
	"context"
	"github.com/olivere/elastic/v7"
)

type AnyESDAO struct {
	client *elastic.Client
}

func (a *AnyESDAO) Input(ctx context.Context, index, docID, data string) error {
	_, err := a.client.Index().Index(index).Id(docID).BodyString(data).Do(ctx)
	return err
}

func NewAnyESDAO(client *elastic.Client) AnyDAO {
	return &AnyESDAO{client: client}
}
