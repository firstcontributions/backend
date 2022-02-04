package mongoqb

import "go.mongodb.org/mongo-driver/bson"

type SearchQuery struct {
	text string
}

func NewSearchQuery(text string) *SearchQuery {
	return &SearchQuery{
		text: text,
	}
}

func (q *SearchQuery) Build() bson.M {
	return bson.M{
		"$text": bson.M{
			"$search": q.text,
		},
	}
}
