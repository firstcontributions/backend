package schema

import (
	"github.com/firstcontributions/backend/internal/models/usersstore"
)

type Reputation struct {
	ref   *usersstore.Reputation
	Value float64
}

func NewReputation(m *usersstore.Reputation) *Reputation {
	if m == nil {
		return nil
	}
	return &Reputation{
		ref:   m,
		Value: m.Value,
	}
}
func (n *Reputation) ToModel() *usersstore.Reputation {
	if n == nil {
		return nil
	}
	return &usersstore.Reputation{
		Value: n.Value,
	}
}
