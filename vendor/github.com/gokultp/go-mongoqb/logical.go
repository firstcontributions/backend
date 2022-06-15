package mongoqb

import "go.mongodb.org/mongo-driver/bson"

const (
	operationAnd = "$and"
	operationOr  = "$or"
	operationNot = "$not"
)

type LogicalQuery struct {
	queries   []IQuery
	operation string
}

func NewLogicalQuery(operation string, queries ...IQuery) *LogicalQuery {
	return &LogicalQuery{
		operation: operation,
		queries:   queries,
	}
}

func (q *LogicalQuery) Build() bson.M {
	if len(q.queries) == 0 {
		return nil
	}
	queries := []bson.M{}
	for _, q := range q.queries {
		if generatedQuery := q.Build(); generatedQuery != nil {
			queries = append(queries, generatedQuery)
		}
	}
	return bson.M{
		q.operation: queries,
	}
}
