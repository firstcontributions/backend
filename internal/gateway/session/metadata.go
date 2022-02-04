package session

import (
	"context"
	"encoding/json"

	"google.golang.org/grpc/metadata"
)

type CxtKey int

const (
	CxtKeySession CxtKey = iota
)

// MetaData encapsulates session info
type MetaData map[string]string

func NewMetaData() MetaData {
	return MetaData{}
}

func FromContext(ctx context.Context) MetaData {
	return ctx.Value(CxtKeySession).(MetaData)
}
func (m MetaData) SetHandle(h string) MetaData {
	m["handle"] = h
	return m
}

func (m MetaData) SetUserID(uid string) MetaData {
	m["user_id"] = uid
	return m
}

func (m MetaData) Handle() string {
	return m["handle"]
}

func (m MetaData) UserID() string {
	return m["user_id"]
}

func (m MetaData) MarshalBinary() ([]byte, error) {
	return json.Marshal(m)
}

func (m *MetaData) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &m)
}

func (m MetaData) Proto() metadata.MD {
	return metadata.New(m)
}
