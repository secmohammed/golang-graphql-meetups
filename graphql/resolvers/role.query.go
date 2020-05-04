package resolvers

import (
    "context"

    "github.com/secmohammed/meetups/graphql"
    "github.com/secmohammed/meetups/graphql/loaders"
    "github.com/secmohammed/meetups/models"
)

type roleResolver struct{ *Resolver }

func (r *Resolver) Role() graphql.RoleResolver {
    return &roleResolver{r}
}
func (r *roleResolver) User(ctx context.Context, obj *models.Role) (*models.User, error) {
    return loaders.GetLoaders(ctx).UserByID.Load(obj.UserID)
}
func (r *roleResolver) Permissions(ctx context.Context, obj *models.Role) (interface{}, error) {
    return obj.Permissions, nil
}
