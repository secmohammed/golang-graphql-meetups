package resolvers

import (
    "context"
    "errors"

    "github.com/secmohammed/meetups/graphql"
    "github.com/secmohammed/meetups/graphql/loaders"
    "github.com/secmohammed/meetups/models"
)

type meetupResolver struct{ *Resolver }

func (m *meetupResolver) User(ctx context.Context, obj *models.Meetup) (*models.User, error) {

    return loaders.GetUserLoader(ctx).Load(obj.UserID)
}

func (r *Resolver) Meetup() graphql.MeetupResolver {
    return &meetupResolver{r}
}

func (r *queryResolver) Meetups(ctx context.Context) ([]*models.Meetup, error) {

    return r.MeetupsRepo.GetMeetups()
}

func (m *mutationResolver) CreateMeetup(ctx context.Context, input graphql.NewMeetup) (*models.Meetup, error) {
    if len(input.Name) < 3 {
        return nil, errors.New("name not long enough")
    }

    if len(input.Description) < 3 {
        return nil, errors.New("description not long enough")
    }

    meetup := &models.Meetup{
        Name:        input.Name,
        Description: input.Description,
        UserID:      "1",
    }

    return m.MeetupsRepo.CreateMeetup(meetup)
}
