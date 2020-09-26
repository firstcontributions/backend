package session

import (
	"encoding/json"
)

// MetaData encapsulates session info
type MetaData struct {
	UserID string `json:"uuid"`
	Handle string `json:"handle"`
}

func (m MetaData) MarshalBinary() ([]byte, error) {
	return json.Marshal(m)
}

func (m *MetaData) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &m)
}
