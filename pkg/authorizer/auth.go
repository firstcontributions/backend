package authorizer

import "github.com/firstcontributions/backend/pkg/graphqlid"

func IsAuthorized(permissions map[string]map[uint8]int, ownership *Scope, entity EnityType, operation int) bool {
	if ownership == nil {
		return true
	}
	if permissions == nil {
		return false
	}

	for _, user := range ownership.Users {
		uid := graphqlid.NewGraphqlID(uint8(UserScope), user, true).String()
		if permissions[uid][uint8(entity)]&operation > 0 {
			return true
		}
	}
	for _, community := range ownership.Communities {
		cid := graphqlid.NewGraphqlID(uint8(CommunityScope), community, true).String()
		if permissions[cid][uint8(entity)]&operation > 0 {
			return true
		}
	}

	return false
}
