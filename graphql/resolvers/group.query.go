package resolvers

import (
    "context"

    "github.com/secmohammed/meetups/graphql"
    "github.com/secmohammed/meetups/graphql/loaders"
    "github.com/secmohammed/meetups/models"
)

type groupResolver struct{ *Resolver }

func (r *Resolver) Group() graphql.GroupResolver {
    return &groupResolver{r}
}
func (r *queryResolver) Groups(ctx context.Context) ([]*models.Group, error) {
    return r.GroupsRepo.GetGroups()
}
func (r *queryResolver) Group(ctx context.Context, id string) (*models.Group, error) {
    return r.GroupsRepo.GetByID(id)
}

func (r *groupResolver) User(ctx context.Context, obj *models.Group) (*models.User, error) {
    return loaders.GetLoaders(ctx).UserByID.Load(obj.UserID)
}
func (r *groupResolver) Categories(ctx context.Context, obj *models.Group) ([]*models.Category, error) {
    return loaders.GetLoaders(ctx).CategoriesByGroup.Load(obj.ID)
}
func (r *groupResolver) Meetups(ctx context.Context, obj *models.Group) ([]*models.Meetup, error) {
    return loaders.GetLoaders(ctx).MeetupsByGroup.Load(obj.ID)
}
func (r *groupResolver) Members(ctx context.Context, obj *models.Group) ([]*models.UserGroup, error) {
    return loaders.GetLoaders(ctx).MembersByGroup.Load(obj.ID)
}
