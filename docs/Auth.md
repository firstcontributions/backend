# authorizer

Ownership
    user: himself
    story: user,community
    comment: created_user, community manager
    community: moderators
    reaction: created user


Multiple type of entities can be the owner
multiple users can own at the same time

Role: {
    name: admin
    entities: [
        {
            name: "*",
            operations: crudm
        }
    ]
}

Role: {
    name: manager
    entities: [
        {
            name: "stories",
            operations: crudm
        },
         {
            name: "comments",
            operations: crudm
        }
    ]
}

User  {
    permissions: [
        {
            role: manager,
            scope: {communities: [c1, c2, c3]}
        },
        {
            role: admin
            scope: {users: [u1]}
        }
    ]
}


Entity : {
    ownership: {
        communities: [c1, c2],
        users: [u1, u2, u3]
    }
}



{
    c1: {stories: 5, comments: 7},
    c2: {stories: 5, comments: 7}
    u1: {stories: 5, comments: 7}
}


c 1
r 1<< 1
u





