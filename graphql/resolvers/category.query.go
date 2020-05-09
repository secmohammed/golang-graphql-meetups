package resolvers

import (
    "context"

    "github.com/secmohammed/meetups/graphql"
    "github.com/secmohammed/meetups/graphql/loaders"
    "github.com/secmohammed/meetups/models"
)

type categoryResolver struct{ *Resolver }

//Category Resolver is used to resolve the propertyies associated to the schema, which are user, and meetups.
func (r *Resolver) Category() graphql.CategoryResolver {
    return &categoryResolver{r}
}

func (c *categoryResolver) Meetups(ctx context.Context, obj *models.Category) ([]*models.Meetup, error) {
    return loaders.GetLoaders(ctx).MeetupsByCategory.Load(obj.ID)
}
func (c *categoryResolver) User(ctx context.Context, obj *models.Category) (*models.User, error) {
    return loaders.GetLoaders(ctx).UserByID.Load(obj.UserID)
}

func (c *queryResolver) Categories(ctx context.Context, limit *int, offset *int) ([]*models.Category, error) {
    return c.CategoriesRepo.GetCategories(limit, offset)
}
func (c *queryResolver) Category(ctx context.Context, name string) (*models.Category, error) {
    return c.CategoriesRepo.GetByName(name)
}
