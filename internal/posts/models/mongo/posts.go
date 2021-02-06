package mongo

import (
	"time"

	"github.com/firstcontributions/backend/internal/posts/proto"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Post defines how post document is stored in mongo
type Post struct {
	UUID        string     `bson:"_id,omitempty"`
	Title       string     `bson:"title,omitempty"`
	Summary     string     `bson:"summary,omitempty"`
	Content     string     `bson:"content,omitempty"`
	CreatedBy   string     `bson:"created_by,omitempty"`
	CoverImage  string     `bson:"cover_image,omitempty"`
	Claps       uint64     `bson:"claps,omitempty"`
	DateUpdated *time.Time `bson:"date_updated,omitempty"`
	DateCreated *time.Time `bson:"date_created,omitempty"`
}

func NewPost() *Post {
	return &Post{}
}

func (p *Post) toProto() *proto.Post {
	return &proto.Post{
		Uuid: p.UUID,
		Data: &proto.PostData{
			Title:      p.Title,
			Summary:    p.Summary,
			Content:    p.Content,
			CoverImage: p.CoverImage,
		},
		CreatedBy: p.CreatedBy,
		Claps:     p.Claps,
		AuditFields: &proto.AuditFields{
			DateCreated: timestamppb.New(*p.DateCreated),
			DateUpdated: timestamppb.New(*p.DateUpdated),
		},
	}
}

func (p *Post) fromProto(post *proto.PostData) error {
	id, err := uuid.NewUUID()
	if err != nil {
		return err
	}
	now := time.Now()
	*p = Post{
		UUID:        id.String(),
		Title:       post.Title,
		Summary:     post.Summary,
		Content:     post.Content,
		CoverImage:  post.CoverImage,
		Claps:       0,
		DateCreated: &now,
	}
	return nil
}
