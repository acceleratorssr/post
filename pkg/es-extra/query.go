package es_extra

import (
	"context"
	"encoding/json"
	"fmt"
	es "github.com/elastic/go-elasticsearch/v8"
	"log"
	"strings"
)

type Query interface {
	Raw() (map[string]any, error)
	Build() (string, error)
}

const ArticleIndexName = "article_index"

// todo 扩展可选参数

type QueryBuilder struct {
	query map[string]interface{}
}

func NewQueryBuilder() *QueryBuilder {
	qb := &QueryBuilder{
		query: make(map[string]interface{}),
	}
	qb.query["query"] = make(map[string]interface{})
	return qb
}

// Match 全文搜索，boost仅第一位有效
func (qb *QueryBuilder) Match(field, value string, boost ...float64) *QueryBuilder {
	matchQuery := map[string]interface{}{
		field: value,
	}

	if len(boost) > 0 && boost[0] > 0 {
		matchQuery["boost"] = boost[0]
	}

	qb.query["query"].(map[string]interface{})["match"] = matchQuery

	return qb
}

// MultiMatch 在一个文档的多个字段（如标题、正文、标签等）中查找关键词
//
//	{
//	  "query": {
//	    "multi_match": {
//	      "query": "keyword",
//	      "fields": [fields...]
//	    }
//	}
func (qb *QueryBuilder) MultiMatch(keyword string, fields []string) *QueryBuilder {
	qb.query["query"].(map[string]interface{})["multi_match"] = map[string]interface{}{
		"query":  keyword,
		"fields": fields,
	}

	return qb
}

func (qb *QueryBuilder) MatchAll() *QueryBuilder {
	qb.query["query"].(map[string]interface{})["match_all"] = map[string]interface{}{
		"match_all": map[string]interface{}{},
	}

	return qb
}

// Term 精确匹配
func (qb *QueryBuilder) Term(field string, value interface{}) *QueryBuilder {
	qb.query["query"].(map[string]interface{})["term"] = map[string]interface{}{
		field: value,
	}

	return qb
}

// Range 范围查询
func (qb *QueryBuilder) Range(field string, gte interface{}, lte interface{}) *QueryBuilder {
	qb.query["query"].(map[string]interface{})["range"] = map[string]interface{}{
		field: map[string]interface{}{
			"gte": gte,
			"lte": lte,
		},
	}

	return qb
}

// Wildcard 模糊匹配
func (qb *QueryBuilder) Wildcard(field, value string) *QueryBuilder {
	qb.query["query"].(map[string]interface{})["wildcard"] = map[string]interface{}{
		field: value,
	}

	return qb
}

// Fuzzy 模糊匹配，容忍拼写错误
// fuzziness 是编辑距离的允许值，可使用 AUTO 根据词的长度动态变化
func (qb *QueryBuilder) Fuzzy(field, value, fuzziness string, boost ...float64) *QueryBuilder {
	fuzzyQuery := map[string]interface{}{
		field: map[string]interface{}{
			"value":     value,
			"fuzziness": fuzziness,
		},
	}

	if len(boost) > 0 && boost[0] > 0 {
		fuzzyQuery[field].(map[string]interface{})["boost"] = boost[0]
	}

	qb.query["query"].(map[string]interface{})["fuzzy"] = fuzzyQuery

	return qb
}

// Exists 查询存在某个字段的文档
func (qb *QueryBuilder) Exists(field string) *QueryBuilder {
	qb.query["query"].(map[string]interface{})["exists"] = map[string]interface{}{
		"field": field,
	}

	return qb
}

// MatchPhrase 查询（短语匹配）
func (qb *QueryBuilder) MatchPhrase(field, value string) *QueryBuilder {
	qb.query["query"].(map[string]interface{})["match_phrase"] = map[string]interface{}{
		field: value,
	}

	return qb
}

// CosineSimilarity 余弦相似度查询，对不同长度的向量（如不同长度的摘要）表现良好
func (qb *QueryBuilder) CosineSimilarity(field string, vector []float32) *QueryBuilder {
	return qb.CustomFunction(field, vector,
		"cosineSimilarity(params.query_vector, 'embedding') + 1.0") // doc[params.field]
}

// DotProduct 计算点积查询，特征向量表示（如图像或文本特征）
func (qb *QueryBuilder) DotProduct(field string, vector []float32) *QueryBuilder {
	return qb.CustomFunction(field, vector,
		"dotProduct(params.query_vector, 'embedding')")
}

// ManhattanDistance 曼哈顿距离查询，特征较稀疏的场景，如推荐系统中的用户特征比较
// 数据较为离散的情况
func (qb *QueryBuilder) ManhattanDistance(field string, vector []float32) *QueryBuilder {
	return qb.CustomFunction(field, vector,
		"l1norm(params.query_vector, 'embedding')")
}

// EuclideanDistance 欧几里得距离查询，聚类分析，数值数据较为均匀且结构化的情况
func (qb *QueryBuilder) EuclideanDistance(field string, vector []float32) *QueryBuilder {
	return qb.CustomFunction(field, vector,
		"l2norm(params.query_vector, 'embedding')")
}

// CustomFunction 自定义计算函数查询
func (qb *QueryBuilder) CustomFunction(field string, vector []float32, customFunction string) *QueryBuilder {
	qb.query["query"].(map[string]interface{})["script_score"] = map[string]interface{}{
		"query": map[string]interface{}{
			"match_all": map[string]interface{}{},
		},
		"script": map[string]interface{}{
			"source": customFunction,
			"params": map[string]interface{}{
				"query_vector": vector,
				"field":        field,
			},
		},
	}

	return qb
}

func (qb *QueryBuilder) KNN(field string, vector []float32, k, numCandidates int) *QueryBuilder {
	delete(qb.query, "query")
	qb.query["knn"] = map[string]interface{}{
		"field":          field,
		"query_vector":   vector,
		"k":              k,
		"num_candidates": numCandidates,
	}

	return qb
}

func (qb *QueryBuilder) Raw() (map[string]any, error) {
	return qb.query["query"].(map[string]any), nil
}

func (qb *QueryBuilder) Build() (string, error) {
	queryJSON, err := json.Marshal(qb.query)
	if err != nil {
		return "", fmt.Errorf("构建查询条件失败: %v", err)
	}
	fmt.Println("生成的查询:", string(queryJSON))
	return string(queryJSON), nil
}

func Search[T any](ctx context.Context, client *es.Client, builder Query, nums int) (T, error) {
	var result T
	query, err := builder.Build()
	if err != nil {
		return result, err
	}

	res, err := client.Search(
		client.Search.WithIndex(ArticleIndexName),
		client.Search.WithBody(strings.NewReader(query)),
		client.Search.WithContext(ctx),
		client.Search.WithSize(nums),
	)
	if err != nil {
		return result, fmt.Errorf("搜索失败: %v", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		log.Printf("search error: %s", res.String())
		return result, fmt.Errorf("搜索失败: %s", res.Status())
	}

	if err = json.NewDecoder(res.Body).Decode(&result); err != nil {
		log.Printf("解析搜索结果失败: %v", err)
	}
	return result, nil
}
