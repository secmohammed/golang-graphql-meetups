package queries

import "github.com/secmohammed/meetups/graphql"

func (r *Resolver) Query() graphql.QueryResolver {
    return &queryResolver{r}
}
