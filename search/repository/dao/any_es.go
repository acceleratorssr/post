package dao

import (
	"context"
	es "github.com/elastic/go-elasticsearch/v8"
)

type AnyESDAO struct {
	client *es.Client
}

func NewAnyESDAO(client *es.Client) AnyDAO {
	return &AnyESDAO{
		client: client,
	}
}

func (a *AnyESDAO) Input(ctx context.Context, index, docId, data string) error {
	//_, err := a.client.Index().
	//	Index(index).Id(docId).BodyString(data).Do(ctx)
	return nil
}
