package mongoqb

import (
	"go.mongodb.org/mongo-driver/bson"
)

const (
	operationEq  = "$eq"
	operationNe  = "$ne"
	operationGt  = "$gt"
	operationGte = "$gte"
	operationLt  = "$lt"
	operationLte = "$lte"
	operationIn  = "$in"
	operationNin = "$nin"
)

type ComparisonQuery struct {
	operation string
	field     string
	operand   interface{}
}

func NewComparisonQuery(
	operation string,
	field string,
	operand interface{},
) *ComparisonQuery {
	return &ComparisonQuery{
		operation: operation,
		field:     field,
		operand:   operand,
	}
}

func (c *ComparisonQuery) Build() bson.M {
	if c.operation == operationEq {
		return bson.M{
			c.field: c.operand,
		}
	}
	return bson.M{
		c.field: bson.M{
			c.operation: c.operand,
		},
	}
}
