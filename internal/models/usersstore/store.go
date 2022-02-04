package usersstore

import "context"

type Store interface {
	// user methods
	CreateUser(context.Context, *User) (*User, error)
	GetUserByID(context.Context, string) (*User, error)
	GetUsers(context.Context, []string, *string, *string, *string, *string, *int64, *int64) ([]*User, bool, bool, string, string, error)
	UpdateUser(context.Context, string, *UserUpdate) error
	DeleteUserByID(context.Context, string) error
	// badge methods
	CreateBadge(context.Context, *Badge) (*Badge, error)
	GetBadgeByID(context.Context, string) (*Badge, error)
	GetBadges(context.Context, []string, *string, *string, *string, *int64, *int64) ([]*Badge, bool, bool, string, string, error)
	DeleteBadgeByID(context.Context, string) error
}
