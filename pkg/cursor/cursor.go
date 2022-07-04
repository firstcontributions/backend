package cursor

import (
	"encoding/json"
)

type Cursor struct {
	Version     int         `json:"v"`
	ID          string      `json:"i"`
	SortBy      string      `json:"s"`
	OffsetValue interface{} `json:"o"`
}

func NewCursor(id, sortBy string, OffsetValue interface{}) *Cursor {
	return &Cursor{
		Version:     1,
		ID:          id,
		SortBy:      sortBy,
		OffsetValue: OffsetValue,
	}
}

func (c *Cursor) String() string {
	bts, _ := json.Marshal(c)
	return string(bts)
}

func FromString(s string) *Cursor {
	c := Cursor{}
	json.Unmarshal([]byte(s), &c)
	return &c
}
