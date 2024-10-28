package dao

import (
	"context"
	_ "embed"
	"fmt"
	es "github.com/elastic/go-elasticsearch/v8"
	"golang.org/x/sync/errgroup"
	es_extra "post/pkg/es-extra"
	"strings"
	"time"
)

var (
	//go:embed article_index.json
	articleIndex string
)

// InitES 创建索引
func InitES(client *es.Client) error {
	const timeout = time.Second * 10
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	var eg errgroup.Group
	eg.Go(func() error {
		return tryCreateIndex(ctx, client, es_extra.ArticleIndexName)
	})

	return eg.Wait()
}

func tryCreateIndex(ctx context.Context, client *es.Client, idxName string) error {
	exists, err := client.Indices.Exists([]string{idxName})
	if err != nil {
		return fmt.Errorf("检测 %s 索引是否存在失败 %w", idxName, err)
	} else if exists.StatusCode == 200 {
		return nil
	} else if exists.StatusCode != 404 {
		return fmt.Errorf("检测 %s 索引是否存在返回错误响应码: %d", idxName, exists.StatusCode)
	}
	defer exists.Body.Close()

	create, err := client.Indices.Create(idxName,
		client.Indices.Create.WithBody(strings.NewReader(articleIndex)),
		client.Indices.Create.WithContext(ctx))
	if err != nil {
		return fmt.Errorf("创建 %s 索引是否存在失败 %w", idxName, err)
	}
	defer create.Body.Close()

	return nil
}
