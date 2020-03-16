package mutations

import (
    "context"
    "errors"

    "github.com/secmohammed/meetups/graphql"
    "github.com/secmohammed/meetups/models"
)

func (m *mutationResolver) CreateMeetup(ctx context.Context, input graphql.NewMeetup) (*models.Meetup, error) {
    if len(input.Name) < 3 {
        return nil, errors.New("Name must be more than 3 characters")
    }
    if len(input.Description) < 3 {
        return nil, errors.New("Description must be more than 3 characters")
    }
    meetup := &models.Meetup{
        Name:        input.Name,
        Description: input.Description,
        UserID:      "1",
    }
    return m.MeetupsRepo.CreateMeetup(meetup)
}
