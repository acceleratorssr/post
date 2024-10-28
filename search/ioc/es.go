package ioc

import (
	es "github.com/elastic/go-elasticsearch/v8"
	"post/search/repository/dao"
)

func InitESClient() *es.Client {
	opts := es.Config{
		Addresses: []string{
			"http://localhost:9200",
		},
		EnableMetrics: true,
		Transport:     nil,
	}
	client, err := es.NewClient(opts)
	if err != nil {
		panic(err)
	}
	dao.InitES(client)

	return client
}
