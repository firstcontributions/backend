package usersstore

import "context"

type Store interface {
	// user methods
	GetUserByID(ctx context.Context, id string) (*User, error)
	GetOneUser(ctx context.Context, filters *UserFilters) (*User, error)
	GetUsers(ctx context.Context, filters *UserFilters, after *string, before *string, first *int64, last *int64, sortBy UserSortBy, sortOrder *string) ([]*User, bool, bool, string, string, error)
	CountUsers(ctx context.Context, filters *UserFilters) (int64, error)
	CreateUser(ctx context.Context, user *User) (*User, error)
	UpdateUser(ctx context.Context, id string, update *UserUpdate) error
	DeleteUserByID(ctx context.Context, id string) error
	// badge methods
	GetBadgeByID(ctx context.Context, id string) (*Badge, error)
	GetOneBadge(ctx context.Context, filters *BadgeFilters) (*Badge, error)
	GetBadges(ctx context.Context, filters *BadgeFilters, after *string, before *string, first *int64, last *int64, sortBy BadgeSortBy, sortOrder *string) ([]*Badge, bool, bool, string, string, error)
	CountBadges(ctx context.Context, filters *BadgeFilters) (int64, error)
	CreateBadge(ctx context.Context, badge *Badge) (*Badge, error)
	UpdateBadge(ctx context.Context, id string, update *BadgeUpdate) error
	DeleteBadgeByID(ctx context.Context, id string) error
}
