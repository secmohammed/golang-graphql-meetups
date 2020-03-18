package resolvers

import (
    "context"

    "github.com/secmohammed/meetups/graphql"
    "github.com/secmohammed/meetups/graphql/loaders"
    "github.com/secmohammed/meetups/models"
)

type commentResolver struct{ *Resolver }

//Meetup resolver.
func (r *Resolver) Meetup() graphql.MeetupResolver {
    return &meetupResolver{r}
}
func (r *Resolver) Comment() graphql.CommentResolver {
    return &commentResolver{r}
}

func (r *queryResolver) Meetups(ctx context.Context, filter *models.MeetupFilter, limit *int, offset *int) ([]*models.Meetup, error) {

    return r.MeetupsRepo.GetMeetups(filter, limit, offset)
}

type meetupResolver struct{ *Resolver }

func (m *meetupResolver) User(ctx context.Context, obj *models.Meetup) (*models.User, error) {
    return loaders.GetLoaders(ctx).UserByID.Load(obj.UserID)
}
func (c *commentResolver) User(ctx context.Context, obj *models.Comment) (*models.User, error) {
    return loaders.GetLoaders(ctx).UserByID.Load(obj.UserID)

}

func (m *meetupResolver) Comments(ctx context.Context, obj *models.Meetup) ([]*models.Comment, error) {
    // loaders.GetLoaders(ctx).CommentsByMeetupID.Load(obj.ID)
    return m.CommentsRepo.GetCommentsForMeetup(obj)
}
