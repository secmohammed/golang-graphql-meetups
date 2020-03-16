//go:generate go run github.com/99designs/gqlgen -v

package resolvers

import (
	"github.com/secmohammed/meetups/graphql"
	"github.com/secmohammed/meetups/postgres"
)

//Resolver struct.
type Resolver struct {
	MeetupsRepo postgres.MeetupsRepo
	UsersRepo   postgres.UsersRepo
}

type mutationResolver struct{ *Resolver }

//
type queryResolver struct{ *Resolver }

func (r *Resolver) Mutation() graphql.MutationResolver {
	return &mutationResolver{r}
}

func (r *Resolver) Query() graphql.QueryResolver {
	return &queryResolver{r}
}
