package userctx

import (
	"context"

	"google.golang.org/grpc/metadata"
)

type Meta struct {
	metadata.MD
}

func (m *Meta) Handle() string {

	return m.getByKey("handle")
}
func (m *Meta) getByKey(key string) string {
	if data := m.Get(key); len(data) == 1 {
		return data[0]
	}
	return ""
}

func (m *Meta) UserID() string {
	return m.getByKey("user_id")
}

func FromIncomingCtx(ctx context.Context) *Meta {
	if m, ok := metadata.FromIncomingContext(ctx); ok {
		return &Meta{
			MD: m,
		}
	}
	return nil
}
