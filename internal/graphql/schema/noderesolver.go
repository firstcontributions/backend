package schema

import (
	"context"
	"errors"

	"github.com/firstcontributions/backend/internal/storemanager"
	"github.com/graph-gophers/graphql-go"
)

type Node interface {
	ID(context.Context) graphql.ID
}

type NodeResolver struct {
	Node
}

type NodeIDInput struct {
	ID graphql.ID
}

func (r *Resolver) Node(ctx context.Context, in NodeIDInput) (*NodeResolver, error) {
	store := storemanager.FromContext(ctx)
	id, err := ParseGraphqlID(in.ID)
	if err != nil {
		return nil, err
	}
	switch id.Type {
	case "badge":
		badgeData, err := store.UsersStore.GetBadgeByID(ctx, id.ID)
		if err != nil {
			return nil, err
		}
		badgeNode := NewBadge(badgeData)
		return &NodeResolver{
			Node: badgeNode,
		}, nil
	case "comment":
		commentData, err := store.StoriesStore.GetCommentByID(ctx, id.ID)
		if err != nil {
			return nil, err
		}
		commentNode := NewComment(commentData)
		return &NodeResolver{
			Node: commentNode,
		}, nil
	case "issue":
		issueData, err := store.IssuesStore.GetIssueByID(ctx, id.ID)
		if err != nil {
			return nil, err
		}
		issueNode := NewIssue(issueData)
		return &NodeResolver{
			Node: issueNode,
		}, nil
	case "reaction":
		reactionData, err := store.StoriesStore.GetReactionByID(ctx, id.ID)
		if err != nil {
			return nil, err
		}
		reactionNode := NewReaction(reactionData)
		return &NodeResolver{
			Node: reactionNode,
		}, nil
	case "story":
		storyData, err := store.StoriesStore.GetStoryByID(ctx, id.ID)
		if err != nil {
			return nil, err
		}
		storyNode := NewStory(storyData)
		return &NodeResolver{
			Node: storyNode,
		}, nil
	case "user":
		userData, err := store.UsersStore.GetUserByID(ctx, id.ID)
		if err != nil {
			return nil, err
		}
		userNode := NewUser(userData)
		return &NodeResolver{
			Node: userNode,
		}, nil
	}
	return nil, errors.New("invalid ID")
}
func (n *NodeResolver) ToBadge() (*Badge, bool) {
	t, ok := n.Node.(*Badge)
	return t, ok
}
func (n *NodeResolver) ToComment() (*Comment, bool) {
	t, ok := n.Node.(*Comment)
	return t, ok
}
func (n *NodeResolver) ToIssue() (*Issue, bool) {
	t, ok := n.Node.(*Issue)
	return t, ok
}
func (n *NodeResolver) ToReaction() (*Reaction, bool) {
	t, ok := n.Node.(*Reaction)
	return t, ok
}
func (n *NodeResolver) ToStory() (*Story, bool) {
	t, ok := n.Node.(*Story)
	return t, ok
}
func (n *NodeResolver) ToUser() (*User, bool) {
	t, ok := n.Node.(*User)
	return t, ok
}
