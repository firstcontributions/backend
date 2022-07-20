package cursor

import (
	"bytes"
	"encoding/base64"
	"math/big"
	"time"

	"github.com/google/uuid"
)

type ValueType uint8

const (
	ValueTypeInt ValueType = iota
	ValueTypeString
	ValueTypeTime
)

func valueToBytes(val interface{}, typ ValueType) []byte {
	switch typ {
	case ValueTypeString:
		return []byte(val.(string))
	case ValueTypeInt:
		return int64ToBytes(val.(int64))
	case ValueTypeTime:
		return int64ToBytes(val.(time.Time).UnixNano())
	}
	return nil
}

func int64ToBytes(v int64) []byte {
	c := new(big.Int)
	c.SetInt64(v)
	return c.Bytes()
}

func bytesToInt64(b []byte) int64 {
	c := new(big.Int)
	c.SetBytes(b)
	return c.Int64()
}

func valueFromBytes(s []byte, typ ValueType) interface{} {
	switch typ {
	case ValueTypeString:
		return string(s)
	case ValueTypeInt:
		return bytesToInt64(s)
	case ValueTypeTime:
		nano := bytesToInt64(s)
		return time.Unix(0, nano)
	}
	return ""
}

type Cursor struct {
	Version     uint8
	ID          string
	SortBy      uint8
	OffsetValue interface{}
	Type        ValueType
}

func NewCursor(id string, sortBy uint8, OffsetValue interface{}, typ ValueType) *Cursor {
	return &Cursor{
		Version:     1,
		ID:          id,
		SortBy:      sortBy,
		OffsetValue: OffsetValue,
		Type:        typ,
	}
}

func (c *Cursor) String() string {
	id, _ := uuid.Parse(c.ID)
	val := []byte{'v', byte(c.Version), '|'}
	idBts, _ := id.MarshalBinary()
	val = append(val, idBts...)
	val = append(val, '|', byte(c.SortBy), '|')
	val = append(val, valueToBytes(c.OffsetValue, c.Type)...)
	val = append(val, '|', byte(c.Type))

	return base64.RawStdEncoding.EncodeToString(val)
}

func FromString(s string) *Cursor {
	if len(s) == 0 {
		return nil
	}
	bts, err := base64.RawStdEncoding.DecodeString(s)
	if err != nil {
		return nil
	}
	parts := bytes.Split(bts, []byte{'|'})

	c := Cursor{}
	c.Version = uint8(parts[0][1])

	id, _ := uuid.FromBytes(parts[1])
	c.ID = id.String()

	c.SortBy = uint8(parts[2][0])
	c.Type = ValueType(parts[4][0])
	c.OffsetValue = valueFromBytes(parts[3], c.Type)
	return &c
}
