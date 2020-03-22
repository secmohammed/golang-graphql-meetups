//go:generate go run github.com/99designs/gqlgen -v

package resolvers

import (
    "sync"

    "github.com/secmohammed/meetups/graphql"
    "github.com/secmohammed/meetups/postgres"
)

//Resolver struct.
type Resolver struct {
    MeetupsRepo       postgres.MeetupsRepo
    UsersRepo         postgres.UsersRepo
    CommentsRepo      postgres.CommentsRepo
    CategoriesRepo    postgres.CategoriesRepo
    AttendeesRepo     postgres.AttendeesRepo
    ConversationsRepo postgres.ConversationsRepo
    mu                sync.Mutex // nolint: structcheck
    Rooms             map[string]*Chatroom
}

type mutationResolver struct{ *Resolver }

//
type queryResolver struct{ *Resolver }

// Mutation method is used to resolve the mutations
func (r *Resolver) Mutation() graphql.MutationResolver {
    return &mutationResolver{r}
}

// Query method is used to resolve the queries
func (r *Resolver) Query() graphql.QueryResolver {
    return &queryResolver{r}
}
