package mutations

import (
    "github.com/secmohammed/meetups/graphql"
    "github.com/secmohammed/meetups/graphql/resolvers"
)

type mutationResolver struct{ *resolvers.Resolver }

func (r *resolvers.Resolver) Mutation() graphql.MutationResolver {
    return &mutationResolver{r}
}
