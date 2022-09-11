package authorizer

const (
	OperationCreate = 1
	OperationRead   = 1 << 1
	OperationUpdate = 1 << 2
	OperaionDelete  = 1 << 3
	OperationManage = 1 << 4
)

type EnityType uint8

const (
	User EnityType = iota
	Story
	Comment
	Badge
	Issue
	Reaction
)

type Role struct {
	Name     string
	Entities []Entity
}

type Entity struct {
	Type       EnityType
	Operations int
}

func NewEntity(typ EnityType, operations string) Entity {
	ops := 0

	for i := range operations {
		switch operations[i] {
		case 'c':
			ops |= OperationCreate
		case 'r':
			ops |= OperationRead
		case 'u':
			ops |= OperationUpdate
		case 'd':
			ops |= OperaionDelete
		case 'm':
			ops |= OperationManage
		}
	}

	return Entity{
		Type:       typ,
		Operations: ops,
	}
}

func GetRole(name string) Role {
	roles := map[string]Role{
		"admin": {
			Name: "admin",
			Entities: []Entity{
				NewEntity(User, "rudm"),
				NewEntity(Badge, "r"),
				NewEntity(Story, "crudm"),
				NewEntity(Comment, "crudm"),
				NewEntity(Issue, "r"),
			},
		},
	}

	return roles[name]
}
