package resolvers

import (
    "context"
    "fmt"

    "github.com/secmohammed/meetups/middlewares"
    "github.com/secmohammed/meetups/models"
    "github.com/secmohammed/meetups/utils/errors"
)

func (m *mutationResolver) DeleteMeetup(ctx context.Context, id string) (bool, error) {
    currentUser, _ := middlewares.GetCurrentUserFromContext(ctx)

    meetup, err := m.MeetupsRepo.GetByID(id)
    if err != nil || meetup == nil {
        return false, errors.ErrRecordNotFound
    }
    if meetup.UserID != currentUser.ID {
        return false, errors.ErrUnauthenticated
    }

    err = m.MeetupsRepo.Delete(meetup)
    if err != nil {
        return false, fmt.Errorf("error while deleting meetup: %v", err)
    }
    return true, nil
}
func (m *mutationResolver) UpdateMeetup(ctx context.Context, id string, input models.UpdateMeetupInput) (*models.Meetup, error) {
    currentUser, _ := middlewares.GetCurrentUserFromContext(ctx)

    meetup, err := m.MeetupsRepo.GetByID(id)

    if err != nil || meetup == nil {
        return nil, errors.ErrRecordNotFound
    }

    if meetup.UserID != currentUser.ID {
        return nil, errors.ErrUnauthenticated
    }
    if err := input.Validate(); err != nil {
        return nil, err
    }
    meetup = &models.Meetup{
        ID:          id,
        Name:        input.Name,
        Description: input.Description,
        StartDate:   *input.StartDate,
        EndDate:     *input.EndDate,
        Location:    input.Location,
    }
    meetup, err = m.MeetupsRepo.Update(meetup)
    if err != nil {
        return nil, errors.ErrInternalError
    }
    return meetup, nil
}
func (m *mutationResolver) CreateMeetup(ctx context.Context, input models.CreateMeetupInput) (*models.Meetup, error) {
    currentUser, _ := middlewares.GetCurrentUserFromContext(ctx)
    if err := input.Validate(); err != nil {
        return nil, err
    }
    meetup := &models.Meetup{
        Name:        input.Name,
        Description: input.Description,
        StartDate:   *input.StartDate,
        EndDate:     *input.EndDate,
        Location:    input.Location,
        UserID:      currentUser.ID,
    }
    meetup, err := m.MeetupsRepo.Create(meetup)
    if input.GroupID != "" {
        group, err := m.GroupsRepo.GetByID(input.GroupID)
        if err != nil {
            return nil, err
        }
        m.GroupsRepo.AttachMeetupToGroup(group, meetup)
    }
    return meetup, err
}
