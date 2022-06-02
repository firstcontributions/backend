package usersstore

import (
	"github.com/firstcontributions/backend/internal/grpc/users/proto"
	"github.com/firstcontributions/backend/internal/models/utils"
)

type Tags struct {
	Languages   []*string `bson:"languages,omitempty"`
	RecentRepos []*string `bson:"recent_repos,omitempty"`
	Topics      []*string `bson:"topics,omitempty"`
}

func NewTags() *Tags {
	return &Tags{}
}
func (tags *Tags) ToProto() *proto.Tags {
	return &proto.Tags{
		Languages:   utils.ToStringArray(tags.Languages),
		RecentRepos: utils.ToStringArray(tags.RecentRepos),
		Topics:      utils.ToStringArray(tags.Topics),
	}
}

func (tags *Tags) FromProto(protoTags *proto.Tags) *Tags {
	tags.Languages = utils.FromStringArray(protoTags.Languages)
	tags.RecentRepos = utils.FromStringArray(protoTags.RecentRepos)
	tags.Topics = utils.FromStringArray(protoTags.Topics)
	return tags
}
