package dao

import (
	"context"
	"encoding/json"
	"github.com/olivere/elastic/v7"
)

type TagESDAO struct {
	client *elastic.Client
}

func NewTagESDAO(client *elastic.Client) TagDAO {
	return &TagESDAO{client: client}
}

func (t *TagESDAO) Search(ctx context.Context, ObjType string, keywords []string) ([]int64, error) {
	query := elastic.NewBoolQuery().Must(
		elastic.NewTermsQueryFromStrings("tags", keywords...),
		elastic.NewTermQuery("obj_type", ObjType))
	resp, err := t.client.Search(TagIndexName).Query(query).Do(ctx)
	if err != nil {
		return nil, err
	}
	res := make([]int64, 0, len(resp.Hits.Hits))
	for _, hit := range resp.Hits.Hits {
		var ele Tags
		err = json.Unmarshal(hit.Source, &ele)
		if err != nil {
			return nil, err
		}
		res = append(res, ele.ObjId)
	}
	return res, nil
}

type Tags struct {
	ObjType string   `json:"obj_type"`
	ObjId   int64    `json:"obj_id"`
	Tags    []string `json:"tags"`
}
