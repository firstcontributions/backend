package mongo

import (
	"context"
	"errors"
	"time"

	"github.com/firstcontributions/backend/internal/profile/proto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Profile maintains the db doc structure
type Profile struct {
	UUID   string `bson:"_id"`
	Name   string `bson:"name"`
	Handle string `bson:"handle"`
	// this field is futuristic, will keep a track of login provider
	//  will need this in later time if we have a conflict for handler
	Provider          string             `bson:"provider"`
	Avatar            string             `bson:"avatar"`
	Reputation        uint64             `bson:"reputation"`
	Badges            []Badge            `bson:"badges"`
	CursorCheckPoints *CursorCheckPoints `bson:"cursor_check_points"`
	DateCreated       time.Time          `bson:"date_created"`
	DateUpdated       time.Time          `bson:"date_updated"`
	Token             *Token             `bson:"token,omitempty"`
}

type CursorCheckPoints struct {
	PullRequest     string `bson:"pull_request"`
	PullRequestFile string `bson:"pull_request_file"`
}

// Badge stores badge db doc structure
type Badge struct {
	UUID       string    `bson:"_id"`
	Name       string    `bson:"name"`
	AssignedOn time.Time `bson:"assigned_on"`
	Progress   uint64    `bson:"progress"`
}

type Token struct {
	AccessToken  string    `bson:"access_token"`
	RefreshToken string    `bson:"refresh_token"`
	TokenType    string    `bson:"token_type"`
	Expiry       time.Time `bson:"expiry"`
}

func convertProfileToModel(p *proto.Profile) *Profile {
	if p == nil {
		return nil
	}
	badges := make([]Badge, len(p.Badges), len(p.Badges))
	for idx, badge := range p.Badges {
		badges[idx] = Badge{
			UUID:       badge.Uuid,
			Name:       badge.Name,
			AssignedOn: badge.AssignedOn.AsTime(),
			Progress:   badge.Progress,
		}
	}
	var token *Token
	if p.GithubToken != nil {
		token = &Token{
			AccessToken:  p.GithubToken.AccessToken,
			RefreshToken: p.GithubToken.RefreshToken,
			TokenType:    p.GithubToken.TokenType,
			Expiry:       p.GithubToken.Expiry.AsTime(),
		}
	}
	var cursor *CursorCheckPoints

	if p.CursorCheckPoints != nil {
		cursor = &CursorCheckPoints{
			PullRequest:     p.CursorCheckPoints.PullRequest,
			PullRequestFile: p.CursorCheckPoints.PullRequestFile,
		}
	}
	return &Profile{
		UUID:              p.Uuid,
		Name:              p.Name,
		Handle:            p.Handle,
		Avatar:            p.Avatar,
		Reputation:        p.Reputation,
		DateCreated:       p.DateCreated.AsTime(),
		DateUpdated:       p.DateUpdated.AsTime(),
		Badges:            badges,
		Token:             token,
		CursorCheckPoints: cursor,
	}
}

// Proto converts Profile model object into proto struct
func (p *Profile) Proto(passSesitiveInformation bool) *proto.Profile {
	if p == nil {
		return nil
	}
	badges := make([]*proto.Badge, len(p.Badges), len(p.Badges))
	for idx, badge := range p.Badges {
		badges[idx] = &proto.Badge{
			Uuid:       badge.UUID,
			Name:       badge.Name,
			AssignedOn: timestamppb.New(badge.AssignedOn),
			Progress:   badge.Progress,
		}
	}
	profile := &proto.Profile{
		Uuid:        p.UUID,
		Name:        p.Name,
		Handle:      p.Handle,
		Avatar:      p.Avatar,
		Reputation:  p.Reputation,
		DateCreated: timestamppb.New(p.DateCreated),
		DateUpdated: timestamppb.New(p.DateUpdated),
		Badges:      badges,
	}
	return profile
}

// GetProfileByUUID gets profile by github handle
func GetProfileByUUID(ctx context.Context, client *mongo.Client, uuid string) (*Profile, error) {
	query := bson.M{
		"_id": uuid,
	}
	var profile Profile
	if err := getCollection(client, CollectionProfile).FindOne(ctx, query).Decode(&profile); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &profile, nil
}

// GetProfileByHandle gets profile by github handle
func GetProfileByHandle(ctx context.Context, client *mongo.Client, handle string) (*Profile, error) {
	query := bson.M{
		"handle": handle,
	}
	var profile Profile
	if err := getCollection(client, CollectionProfile).FindOne(ctx, query).Decode(&profile); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &profile, nil
}

// CreateProfile profile creates a profile
func CreateProfile(ctx context.Context, client *mongo.Client, profile *proto.Profile) error {
	mProfile := convertProfileToModel(profile)
	mProfile.DateCreated = time.Now()
	mProfile.DateUpdated = time.Now()
	_, err := getCollection(client, CollectionProfile).InsertOne(ctx, mProfile)
	return err
}

func UpdateProfile(ctx context.Context, client *mongo.Client, profile *proto.Profile) error {
	mProfile := convertProfileToModel(profile)
	mProfile.DateUpdated = time.Now()
	update := map[string]*Profile{
		"$set": mProfile,
	}
	query := map[string]string{
		"_id": mProfile.UUID,
	}
	_, err := getCollection(client, CollectionProfile).UpdateOne(ctx, query, update, &options.UpdateOptions{})
	return err
}
