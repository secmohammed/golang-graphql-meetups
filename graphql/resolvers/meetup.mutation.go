package resolvers

import (
    "context"
    "errors"
    "fmt"

    "github.com/secmohammed/meetups/middlewares"
    "github.com/secmohammed/meetups/models"
)

var (
    // ErrUnauthenticated is used to indicate that user is unauthenticated.
    ErrUnauthenticated = errors.New("Unauthorized Attempt")
)

func (m *mutationResolver) DeleteMeetup(ctx context.Context, id string) (bool, error) {
    currentUser, err := middlewares.GetCurrentUserFromContext(ctx)
    if err != nil {
        return false, ErrUnauthenticated
    }

    meetup, err := m.MeetupsRepo.GetByID(id)
    if err != nil || meetup == nil {
        return false, errors.New("meetup doesn't exist")
    }
    if meetup.UserID != currentUser.ID {
        return false, errors.New("Unauthorized attempt")
    }

    err = m.MeetupsRepo.DeleteMeetup(meetup)
    if err != nil {
        return false, fmt.Errorf("error while deleting meetup: %v", err)
    }
    return true, nil
}
func (m *mutationResolver) UpdateMeetup(ctx context.Context, id string, input models.UpdateMeetup) (*models.Meetup, error) {
    currentUser, err := middlewares.GetCurrentUserFromContext(ctx)
    if err != nil {
        return nil, ErrUnauthenticated
    }

    meetup, err := m.MeetupsRepo.GetByID(id)

    if err != nil || meetup == nil {
        return nil, errors.New("meetup doesn't exist")
    }

    if meetup.UserID != currentUser.ID {
        return nil, errors.New("Unauthorized attempt")
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
func (m *mutationResolver) CreateMeetup(ctx context.Context, input models.NewMeetup) (*models.Meetup, error) {
    currentUser, err := middlewares.GetCurrentUserFromContext(ctx)
    if err != nil {
        return nil, ErrUnauthenticated
    }

    meetup := &models.Meetup{
        Name:        input.Name,
        Description: input.Description,
        UserID:      currentUser.ID,
    }
    if err := meetup.Validate(); err != nil {
        return nil, err
    }

    return m.MeetupsRepo.CreateMeetup(meetup)
}
