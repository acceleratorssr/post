package es_extra

import (
	"encoding/json"
	"fmt"
)

type BoolQuery struct {
	mustClauses        []Query
	shouldClauses      []Query
	mustNotClauses     []Query
	filterClauses      []Query
	query              map[string]interface{}
	boost              *float64
	minimumShouldMatch *int
	queryName          string
}

func NewBoolQuery() *BoolQuery {
	return &BoolQuery{
		query:          make(map[string]interface{}),
		mustClauses:    make([]Query, 0),
		shouldClauses:  make([]Query, 0),
		mustNotClauses: make([]Query, 0),
		filterClauses:  make([]Query, 0),
	}
}

func (q *BoolQuery) Must(queries ...Query) *BoolQuery {
	q.mustClauses = append(q.mustClauses, queries...)
	return q
}

func (q *BoolQuery) Should(queries ...Query) *BoolQuery {
	q.shouldClauses = append(q.shouldClauses, queries...)
	return q
}

func (q *BoolQuery) MustNot(queries ...Query) *BoolQuery {
	q.mustNotClauses = append(q.mustNotClauses, queries...)
	return q
}

func (q *BoolQuery) Filter(queries ...Query) *BoolQuery {
	q.filterClauses = append(q.filterClauses, queries...)
	return q
}

func (q *BoolQuery) Boost(boost float64) *BoolQuery {
	q.boost = &boost
	return q
}

func (q *BoolQuery) MinimumShouldMatch(minimum int) *BoolQuery {
	q.minimumShouldMatch = &minimum
	return q
}

func (q *BoolQuery) QueryName(name string) *BoolQuery {
	q.queryName = name
	return q
}

// Raw 生成 BoolQuery 的查询结构
func (q *BoolQuery) Raw() (map[string]any, error) {
	boolQuery := make(map[string]interface{})

	// must
	if len(q.mustClauses) > 0 {
		var mustQueries []interface{}
		for _, clause := range q.mustClauses {
			src, err := clause.Raw()
			if err != nil {
				return nil, err
			}
			mustQueries = append(mustQueries, src)
		}
		boolQuery["must"] = mustQueries
	}

	// should
	if len(q.shouldClauses) > 0 {
		var shouldQueries []interface{}
		for _, clause := range q.shouldClauses {
			src, err := clause.Raw()
			if err != nil {
				return nil, err
			}
			shouldQueries = append(shouldQueries, src)
		}
		boolQuery["should"] = shouldQueries
	}

	// must_not
	if len(q.mustNotClauses) > 0 {
		var mustNotQueries []interface{}
		for _, clause := range q.mustNotClauses {
			src, err := clause.Raw()
			if err != nil {
				return nil, err
			}
			mustNotQueries = append(mustNotQueries, src)
		}
		boolQuery["must_not"] = mustNotQueries
	}

	// filter
	if len(q.filterClauses) > 0 {
		var filterQueries []interface{}
		for _, clause := range q.filterClauses {
			src, err := clause.Raw()
			if err != nil {
				return nil, err
			}
			filterQueries = append(filterQueries, src)
		}
		boolQuery["filter"] = filterQueries
	}

	if q.boost != nil {
		boolQuery["boost"] = *q.boost
	}
	if q.minimumShouldMatch != nil {
		boolQuery["minimum_should_match"] = *q.minimumShouldMatch
	}
	if q.queryName != "" {
		boolQuery["_name"] = q.queryName
	}
	q.query = map[string]interface{}{"bool": boolQuery}
	q.query = map[string]interface{}{"query": q.query}

	/*
		{
		  "query" : {
		    "bool" : {
		      "should" : [ {
		        "multi_match" : {
		          "fields" : [ "title", "content" ],
		          "query" : "backend"
		        }
		      }, {
		        "fuzzy" : {
		          "title" : {
		            "fuzziness" : "2",
		            "value" : "haalth"
		          }
		        }
		      } ]
		    }
		  }
		}
	*/
	return q.query, nil
}

func (q *BoolQuery) Build() (string, error) {
	if q.query == nil {
		q.query, _ = q.Raw()
	}

	queryJSON, err := json.Marshal(q.query)

	if err != nil {
		return "", fmt.Errorf("构建查询条件失败: %v", err)
	}
	return string(queryJSON), nil
}
