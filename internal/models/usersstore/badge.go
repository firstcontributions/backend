package usersstore

import (
	"time"

	"github.com/firstcontributions/backend/internal/grpc/users/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Badge struct {
	UserID                        string    `bson:"user_id"`
	CurrentLevel                  int64     `bson:"current_level,omitempty"`
	DisplayName                   string    `bson:"display_name,omitempty"`
	Id                            string    `bson:"_id"`
	Points                        int64     `bson:"points,omitempty"`
	ProgressPercentageToNextLevel int64     `bson:"progress_percentage_to_next_level,omitempty"`
	TimeCreated                   time.Time `bson:"time_created,omitempty"`
	TimeUpdated                   time.Time `bson:"time_updated,omitempty"`
}

func NewBadge() *Badge {
	return &Badge{}
}
func (badge *Badge) ToProto() *proto.Badge {
	return &proto.Badge{
		UserId:                        badge.UserID,
		CurrentLevel:                  badge.CurrentLevel,
		DisplayName:                   badge.DisplayName,
		Id:                            badge.Id,
		Points:                        badge.Points,
		ProgressPercentageToNextLevel: badge.ProgressPercentageToNextLevel,
		TimeCreated:                   timestamppb.New(badge.TimeCreated),
		TimeUpdated:                   timestamppb.New(badge.TimeUpdated),
	}
}

func (badge *Badge) FromProto(protoBadge *proto.Badge) *Badge {
	badge.UserID = protoBadge.UserId
	badge.CurrentLevel = protoBadge.CurrentLevel
	badge.DisplayName = protoBadge.DisplayName
	badge.Id = protoBadge.Id
	badge.Points = protoBadge.Points
	badge.ProgressPercentageToNextLevel = protoBadge.ProgressPercentageToNextLevel
	badge.TimeCreated = protoBadge.TimeCreated.AsTime()
	badge.TimeUpdated = protoBadge.TimeUpdated.AsTime()
	return badge
}

type BadgeUpdate struct {
	CurrentLevel                  *int64     `bson:"current_level,omitempty"`
	Points                        *int64     `bson:"points,omitempty"`
	ProgressPercentageToNextLevel *int64     `bson:"progress_percentage_to_next_level,omitempty"`
	TimeUpdated                   *time.Time `bson:"time_updated,omitempty"`
}

func (badge *BadgeUpdate) ToProto() *proto.UpdateBadgeRequest {
	p := &proto.UpdateBadgeRequest{}
	if badge.CurrentLevel != nil {
		p.CurrentLevel = badge.CurrentLevel
	}
	if badge.Points != nil {
		p.Points = badge.Points
	}
	if badge.ProgressPercentageToNextLevel != nil {
		p.ProgressPercentageToNextLevel = badge.ProgressPercentageToNextLevel
	}
	if badge.TimeUpdated != nil {
		p.TimeUpdated = timestamppb.New(*badge.TimeUpdated)
	}
	return p
}

func (badge *BadgeUpdate) FromProto(protoBadge *proto.UpdateBadgeRequest) {
	badge.CurrentLevel = protoBadge.CurrentLevel
	badge.Points = protoBadge.Points
	badge.ProgressPercentageToNextLevel = protoBadge.ProgressPercentageToNextLevel
	time_updated := protoBadge.TimeUpdated.AsTime()
	badge.TimeUpdated = &time_updated
}
