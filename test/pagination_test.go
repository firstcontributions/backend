package test

import (
	"context"
	"os"
	"strconv"
	"testing"

	"github.com/firstcontributions/backend/internal/models/usersstore"
	usersmongo "github.com/firstcontributions/backend/internal/models/usersstore/mongo"
	"github.com/firstcontributions/backend/pkg/cursor"
	"github.com/google/uuid"

	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PaginationTestSuite struct {
	suite.Suite
	badges []*usersstore.Badge
	usersstore.Store
	client *mongo.Client
}

func (t *PaginationTestSuite) SetupTest() {
	ctx := context.Background()
	client, err := mongo.NewClient(options.Client().ApplyURI(os.Getenv("MONGO_URL")))
	t.NoError(err)
	t.NoError(client.Connect(ctx))
	t.createBadges(ctx, client)
	t.Store, err = usersmongo.NewUsersStore(ctx, os.Getenv("MONGO_URL"))
	t.NoError(err)
	t.client = client
}

func (t *PaginationTestSuite) TearDownTest() {
	t.client.Database("users").Drop(context.TODO())
}

func (t *PaginationTestSuite) createBadges(ctx context.Context, client *mongo.Client) {
	for i := 1; i <= 10; i++ {
		id := strconv.Itoa(i)
		badge := &usersstore.Badge{
			Id:          uuid.NewString(),
			DisplayName: id,
			Points:      int64(i),
		}
		t.badges = append(t.badges, badge)
		_, err := client.Database("users").Collection("badges").InsertOne(ctx, badge)
		t.NoError(err)
	}
}

func (t *PaginationTestSuite) TestSortBy() {
	three := int64(3)
	tests := []struct {
		name      string
		after     *cursor.Cursor
		before    *cursor.Cursor
		sortOrder string
		sortBy    usersstore.BadgeSortBy
		first     *int64
		last      *int64
		want      []string
	}{
		{
			name:      "asc order of points at",
			sortOrder: "asc",
			sortBy:    usersstore.BadgeSortByPoints,
			first:     &three,
			want:      []string{"1", "2", "3"},
		},
		{
			name:      "asc order of points at with after",
			sortOrder: "asc",
			sortBy:    usersstore.BadgeSortByPoints,
			first:     &three,
			after:     cursor.NewCursor(t.badges[2].Id, uint8(usersstore.BadgeSortByPoints), t.badges[2].Points, cursor.ValueTypeInt),
			want:      []string{"4", "5", "6"},
		},
		{
			name:      "asc order of points at with before",
			sortOrder: "asc",
			sortBy:    usersstore.BadgeSortByPoints,
			last:      &three,
			want:      []string{"8", "9", "10"},
		},
		{
			name:      "asc order of points at with before and last",
			sortOrder: "asc",
			sortBy:    usersstore.BadgeSortByPoints,
			last:      &three,
			before:    cursor.NewCursor(t.badges[7].Id, uint8(usersstore.BadgeSortByPoints), t.badges[7].Points, cursor.ValueTypeInt),
			want:      []string{"5", "6", "7"},
		},

		// desc order cases
		{
			name:      "desc order of points at",
			sortOrder: "desc",
			sortBy:    usersstore.BadgeSortByPoints,
			first:     &three,
			want:      []string{"10", "9", "8"},
		},
		{
			name:      "desc order of created at with after",
			sortOrder: "desc",
			sortBy:    usersstore.BadgeSortByPoints,
			first:     &three,
			after:     cursor.NewCursor(t.badges[7].Id, uint8(usersstore.BadgeSortByPoints), t.badges[7].Points, cursor.ValueTypeInt),
			want:      []string{"7", "6", "5"},
		},
		{
			name:      "desc order of created at with last",
			sortOrder: "desc",
			sortBy:    usersstore.BadgeSortByPoints,
			last:      &three,
			want:      []string{"3", "2", "1"},
		},
		{
			name:      "desc order of created at with before and last",
			sortOrder: "desc",
			sortBy:    usersstore.BadgeSortByPoints,
			last:      &three,
			before:    cursor.NewCursor(t.badges[2].Id, uint8(usersstore.BadgeSortByPoints), t.badges[2].Points, cursor.ValueTypeInt),
			want:      []string{"6", "5", "4"},
		},
	}

	for _, tcase := range tests {
		var after, before *string
		if tcase.after != nil {
			tmp := tcase.after.String()
			after = &tmp
		}
		if tcase.before != nil {
			tmp := tcase.before.String()
			before = &tmp
		}
		badges, _, _, _, err := t.Store.GetBadges(
			context.Background(),
			&usersstore.BadgeFilters{},
			after,
			before,
			tcase.first,
			tcase.last,
			usersstore.BadgeSortByPoints,
			&tcase.sortOrder,
		)
		t.NoError(err)
		t.NoError(err)
		resultIds := []string{}
		for _, badge := range badges {
			resultIds = append(resultIds, badge.DisplayName)
		}
		t.Assert().Equal(tcase.want, resultIds, tcase.name)
	}

}

func TestPaginationTestSuite(t *testing.T) {
	suite.Run(t, new(PaginationTestSuite))
}
