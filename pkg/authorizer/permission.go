package authorizer

import "github.com/firstcontributions/backend/pkg/graphqlid"

type Permission struct {
	Role  string
	Scope Scope
}

type Scope struct {
	Users       []string
	Communities []string
}

type AuthScope uint8

const (
	UserScope AuthScope = iota
	CommunityScope
)

func GetResolvedUserPermission(permissions []Permission) map[string]map[uint8]int {
	res := map[string]map[uint8]int{}
	for _, p := range permissions {
		role := GetRole(p.Role)
		for _, e := range role.Entities {
			for _, u := range p.Scope.Users {
				uid := graphqlid.NewGraphqlID(uint8(UserScope), u, true).String()
				if res[uid] == nil {
					res[uid] = map[uint8]int{}
				}
				res[uid][uint8(e.Type)] |= e.Operations
			}
			for _, u := range p.Scope.Communities {
				cid := graphqlid.NewGraphqlID(uint8(CommunityScope), u, true).String()
				if res[cid] == nil {
					res[cid] = map[uint8]int{}
				}
				res[cid][uint8(e.Type)] |= e.Operations
			}
		}
	}
	return res
}
