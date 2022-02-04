package mongoqb

import "go.mongodb.org/mongo-driver/bson"

type IQuery interface {
	Build() bson.M
}

type QueryBuilder struct {
	queries []IQuery
}

func NewQueryBuilder() *QueryBuilder {
	return &QueryBuilder{}
}

func (q *QueryBuilder) Build() bson.M {
	return NewLogicalQuery(operationAnd, q.queries...).Build()
}

func (q *QueryBuilder) Queries() []IQuery {
	return q.queries
}

func (q *QueryBuilder) Eq(field string, operand interface{}) *QueryBuilder {
	q.queries = append(
		q.queries,
		NewComparisonQuery(operationEq, field, operand),
	)
	return q
}

func (q *QueryBuilder) Ne(field string, operand interface{}) *QueryBuilder {
	q.queries = append(
		q.queries,
		NewComparisonQuery(operationNe, field, operand),
	)
	return q
}

func (q *QueryBuilder) Lt(field string, operand interface{}) *QueryBuilder {
	q.queries = append(
		q.queries,
		NewComparisonQuery(operationLt, field, operand),
	)
	return q
}

func (q *QueryBuilder) Lte(field string, operand interface{}) *QueryBuilder {
	q.queries = append(
		q.queries,
		NewComparisonQuery(operationLte, field, operand),
	)
	return q
}

func (q *QueryBuilder) Gt(field string, operand interface{}) *QueryBuilder {
	q.queries = append(
		q.queries,
		NewComparisonQuery(operationGt, field, operand),
	)
	return q
}

func (q *QueryBuilder) Gte(field string, operand interface{}) *QueryBuilder {
	q.queries = append(
		q.queries,
		NewComparisonQuery(operationGte, field, operand),
	)
	return q
}

func (q *QueryBuilder) In(field string, operand interface{}) *QueryBuilder {
	q.queries = append(
		q.queries,
		NewComparisonQuery(operationIn, field, operand),
	)
	return q
}

func (q *QueryBuilder) Nin(field string, operand interface{}) *QueryBuilder {
	q.queries = append(
		q.queries,
		NewComparisonQuery(operationNin, field, operand),
	)
	return q
}

func (q *QueryBuilder) And(queries ...IQuery) *QueryBuilder {
	q.queries = append(
		q.queries,
		NewLogicalQuery(operationAnd, queries...),
	)
	return q
}

func (q *QueryBuilder) Or(queries ...IQuery) *QueryBuilder {
	q.queries = append(
		q.queries,
		NewLogicalQuery(operationOr, queries...),
	)
	return q
}

func (q *QueryBuilder) Not(queries ...IQuery) *QueryBuilder {
	q.queries = append(
		q.queries,
		NewLogicalQuery(operationNot, queries...),
	)
	return q
}

// Make sure text index is created for the fields to be queried
// eg: db.articles.createIndex( { subject: "text" , author: "text" } )
func (q *QueryBuilder) Search(value string) *QueryBuilder {
	q.queries = append(q.queries, NewSearchQuery(value))
	return q
}
