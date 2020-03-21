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

    err = m.MeetupsRepo.Delete(meetup)
    if err != nil {
        return false, fmt.Errorf("error while deleting meetup: %v", err)
    }
    return true, nil
}
func (m *mutationResolver) UpdateMeetup(ctx context.Context, id string, input models.UpdateMeetupInput) (*models.Meetup, error) {
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
    if err := input.Validate(); err != nil {
        return nil, err
    }
    meetup = &models.Meetup{
        ID:          id,
        Name:        input.Name,
        Description: input.Description,
        StartDate:   input.StartDate,
        EndDate:     input.EndDate,
        Location:    input.Location,
    }
    meetup, err = m.MeetupsRepo.Update(meetup)
    if err != nil {
        return nil, fmt.Errorf("error while updating meetup: %v", err)
    }
    return meetup, nil
}
func (m *mutationResolver) CreateMeetup(ctx context.Context, input models.CreateMeetupInput) (*models.Meetup, error) {
    currentUser, err := middlewares.GetCurrentUserFromContext(ctx)
    if err != nil {
        return nil, ErrUnauthenticated
    }
    if err := input.Validate(); err != nil {
        return nil, err
    }
    meetup := &models.Meetup{
        Name:        input.Name,
        Description: input.Description,
        StartDate:   input.StartDate,
        EndDate:     input.EndDate,
        Location:    input.Location,
        UserID:      currentUser.ID,
    }

    return m.MeetupsRepo.Create(meetup)
}
