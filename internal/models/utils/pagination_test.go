package utils

import (
	"fmt"
	"sort"
	"testing"

	"github.com/firstcontributions/backend/pkg/cursor"
	"github.com/gokultp/go-mongoqb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func queryData(query mongoqb.QueryBuilder, opt options.FindOptions) []map[string]string {
	data := []map[string]string{
		{
			"_id":          "1",
			"time_created": "2022-01-01",
			"time_updated": "2022-06-01",
		},
		{
			"_id":          "2",
			"time_created": "2022-01-02",
			"time_updated": "2022-05-29",
		},
		{
			"_id":          "3",
			"time_created": "2022-01-03",
			"time_updated": "2022-05-28",
		},
		{
			"_id":          "4",
			"time_created": "2022-01-04",
			"time_updated": "2022-05-27",
		},
		{
			"_id":          "5",
			"time_created": "2022-01-05",
			"time_updated": "2022-05-26",
		},
		{
			"_id":          "6",
			"time_created": "2022-01-06",
			"time_updated": "2022-05-25",
		},
		{
			"_id":          "7",
			"time_created": "2022-01-07",
			"time_updated": "2022-05-24",
		},
		{
			"_id":          "8",
			"time_created": "2022-01-08",
			"time_updated": "2022-05-23",
		},
	}
	sort.Slice(data, func(i, j int) bool {
		for _, s := range opt.Sort.(bson.D) {
			if s.Value.(int) > 0 {
				if data[i][s.Key] < data[j][s.Key] {
					return true
				}
				if data[i][s.Key] > data[j][s.Key] {
					return false
				}
			} else {
				if data[i][s.Key] > data[j][s.Key] {
					return true
				}
				if data[i][s.Key] < data[j][s.Key] {
					return false
				}
			}
		}
		return false
	})

	for _, q := range query.Queries() {
		cq := q.(*mongoqb.ComparisonQuery)
		tmp := []map[string]string{}
		for field, cond := range cq.Build() {
			for op, operand := range cond.(bson.M) {
				for _, d := range data {
					val := d[field]
					operandVal := operand.(string)
					if op == "$gte" && val >= operandVal {
						tmp = append(tmp, d)
					} else if op == "$lte" && val <= operandVal {
						tmp = append(tmp, d)
					}
				}

			}
		}
		data = tmp
	}
	return data[:min(int(*opt.Limit), len(data))]

}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func getData(after, before, sortOrder, sortBy *string, first, last *int64) ([]map[string]string, bool, bool, string, string, error) {
	qb := mongoqb.NewQueryBuilder()
	limit, order, cursorStr := GetLimitAndSortOrderAndCursor(first, last, after, before)
	var c *cursor.Cursor
	if cursorStr != nil {
		c = cursor.FromString(*cursorStr)
		if c != nil {
			if order == 1 {
				qb.Gte(c.SortBy, c.OffsetValue)
			} else {
				qb.Lte(c.SortBy, c.OffsetValue)
			}
		}
	}
	limit += 2
	options := &options.FindOptions{
		Limit: &limit,
		Sort:  GetSortOrder(sortBy, sortOrder, order),
	}
	fmt.Println(options, qb)

	var firstCursor, lastCursor string
	var hasNextPage, hasPreviousPage bool

	results := queryData(*qb, *options)
	count := len(results)
	if count == 0 {
		return results, hasNextPage, hasPreviousPage, firstCursor, lastCursor, nil
	}

	// check if the cursor element present, if yes that can be a prev elem
	if c != nil && results[0]["_id"] == c.ID {
		hasPreviousPage = true
		results = results[1:]
		count--
	} else if c != nil && results[count-1]["_id"] == c.ID {
		hasNextPage = true
		results = results[:count-1]
		count--
	}

	// check if actual limit +1 elements are there, if yes trim it to limit
	if count >= int(limit)-1 {
		hasNextPage = true
		results = results[:limit-2]
		count = len(results)
	}

	if count > 0 {
		firstCursor = cursor.NewCursor(results[0]["_id"], *sortBy, results[0][*sortBy]).String()
		lastCursor = cursor.NewCursor(results[count-1]["_id"], *sortBy, results[count-1][*sortBy]).String()
	}
	if order < 0 {
		hasNextPage, hasPreviousPage = hasPreviousPage, hasNextPage
		firstCursor, lastCursor = lastCursor, firstCursor
		results = reverseMapList(results)
	}
	return results, hasNextPage, hasPreviousPage, firstCursor, lastCursor, nil
}

func reverseMapList(list []map[string]string) []map[string]string {
	ln := len(list)
	for i := 0; i < ln/2; i++ {
		list[i], list[ln-i-1] = list[ln-i-1], list[i]
	}
	return list
}

func Test_SortOrder_Mock(t *testing.T) {
	cursorTimeCreated2 := cursor.NewCursor("2", "time_created", "2022-01-02").String()
	cursorTimeUpdated3 := cursor.NewCursor("3", "time_updated", "2022-05-28").String()

	cursorTimeCreated5 := cursor.NewCursor("5", "time_created", "2022-01-05").String()
	cursorTimeUpdated5 := cursor.NewCursor("5", "time_updated", "2022-05-26").String()

	asc := "asc"
	timeCreated := "time_created"
	timeUpdated := "time_updated"

	first := int64(3)
	last := int64(2)

	desc := "desc"

	tests := []struct {
		name      string
		after     *string
		before    *string
		sortOrder *string
		sortBy    *string
		first     *int64
		last      *int64
		want      []string
	}{
		{
			name:      "should work for asc order of time_created with after cursor",
			after:     &cursorTimeCreated2,
			sortOrder: &asc,
			sortBy:    &timeCreated,
			first:     &first,
			want:      []string{"3", "4", "5"},
		},
		{
			name:      "should work for desc order of time_created with after cursor",
			after:     &cursorTimeCreated2,
			sortOrder: &desc,
			sortBy:    &timeCreated,
			first:     &first,
			want:      []string{"8", "7", "6"},
		},
		{
			name:      "should work for asc order of time_updated with after cursor",
			after:     &cursorTimeUpdated3,
			sortOrder: &asc,
			sortBy:    &timeUpdated,
			first:     &first,
			want:      []string{"2", "1"},
		},
		{
			name:      "should work for desc order of time_updated with after cursor",
			after:     &cursorTimeUpdated3,
			sortOrder: &desc,
			sortBy:    &timeUpdated,
			first:     &first,
			want:      []string{"1", "2"},
		},
		{
			name:      "should work for asc order of time_created with before cursor",
			before:    &cursorTimeCreated5,
			sortOrder: &asc,
			sortBy:    &timeCreated,
			last:      &last,
			want:      []string{"3", "4"},
		},
		{
			name:      "should work for desc order of time_created with before cursor",
			before:    &cursorTimeCreated5,
			sortOrder: &desc,
			sortBy:    &timeCreated,
			last:      &last,
			want:      []string{"2", "1"},
		},
		{
			name:      "should work for asc order of time_updated with before cursor",
			before:    &cursorTimeUpdated5,
			sortOrder: &asc,
			sortBy:    &timeUpdated,
			last:      &last,
			want:      []string{"7", "6"},
		},
		{
			name:      "should work for desc order of time_updated with before cursor",
			before:    &cursorTimeUpdated5,
			sortOrder: &desc,
			sortBy:    &timeUpdated,
			last:      &last,
			want:      []string{"7", "8"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, _, _, _, _, _ := getData(tt.after, tt.before, tt.sortOrder, tt.sortBy, tt.first, tt.last)
			if len(data) != len(tt.want) {
				t.Error("length of result did not match")
				return
			}
			for i := 0; i < len(tt.want); i++ {
				if data[i]["_id"] != tt.want[i] {
					t.Errorf("expected %v but got %v", tt.want[i], data[i]["_id"])
				}
			}

		})
	}

}

func Test_getSortOrderFromString(t *testing.T) {
	desc := "desc"
	any := "any"
	tests := []struct {
		name  string
		order *string
		want  int
	}{
		{
			name: "should be ascending order if order is null",
			want: 1,
		},
		{
			name:  "should be descending order if order is desc",
			order: &desc,
			want:  -1,
		},
		{
			name:  "should be ascending order by default",
			order: &any,
			want:  1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getSortOrderFromString(tt.order); got != tt.want {
				t.Errorf("getSortOrderFromString() = %v, want %v", got, tt.want)
			}
		})
	}
}
