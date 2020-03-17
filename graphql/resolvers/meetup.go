package resolvers

import (
    "context"
    "errors"
    "fmt"

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

func (r *queryResolver) Meetups(ctx context.Context, filter *graphql.MeetupFilter, limit *int, offset *int) ([]*models.Meetup, error) {

    return r.MeetupsRepo.GetMeetups(filter, limit, offset)
}
func (m *mutationResolver) DeleteMeetup(ctx context.Context, id string) (bool, error) {
    meetup, err := m.MeetupsRepo.GetByID(id)
    if err != nil || meetup == nil {
        return false, errors.New("meetup doesn't exist")
    }
    err = m.MeetupsRepo.DeleteMeetup(meetup)
    if err != nil {
        return false, fmt.Errorf("error while deleting meetup: %v", err)
    }
    return true, nil
}
func (m *mutationResolver) UpdateMeetup(ctx context.Context, id string, input graphql.UpdateMeetup) (*models.Meetup, error) {
    meetup, err := m.MeetupsRepo.GetByID(id)
    if err != nil || meetup == nil {
        return nil, errors.New("meetup doesn't exist")
    }

    meetup.Name = input.Name
    meetup.Description = input.Description
    if err := meetup.Validate(); err != nil {
        return nil, err
    }
    meetup, err = m.MeetupsRepo.Update(meetup)
    if err != nil {
        return nil, fmt.Errorf("error while updating meetup: %v", err)
    }
    return meetup, nil
}
func (m *mutationResolver) CreateMeetup(ctx context.Context, input graphql.NewMeetup) (*models.Meetup, error) {

    meetup := &models.Meetup{
        Name:        input.Name,
        Description: input.Description,
        UserID:      "1",
    }
    if err := meetup.Validate(); err != nil {
        return nil, err
    }

    return m.MeetupsRepo.CreateMeetup(meetup)
}
