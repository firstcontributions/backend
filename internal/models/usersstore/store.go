package usersstore

import "context"

type Store interface {
	// user methods
	GetUserByID(ctx context.Context, id string) (*User, error)
	GetUsers(ctx context.Context, ids []string, search *string, handle *string, after *string, before *string, first *int64, last *int64) ([]*User, bool, bool, string, string, error)
	CreateUser(ctx context.Context, user *User) (*User, error)
	UpdateUser(ctx context.Context, id string, update *UserUpdate) error
	DeleteUserByID(ctx context.Context, id string) error
	// badge methods
	GetBadgeByID(ctx context.Context, id string) (*Badge, error)
	GetBadges(ctx context.Context, ids []string,
		user *User, after *string, before *string, first *int64, last *int64) ([]*Badge, bool, bool, string, string, error)
	CreateBadge(ctx context.Context, badge *Badge) (*Badge, error)
	UpdateBadge(ctx context.Context, id string, update *BadgeUpdate) error
	DeleteBadgeByID(ctx context.Context, id string) error
}
