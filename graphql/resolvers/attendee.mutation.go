package resolvers

import (
    "context"
    "errors"

    "github.com/secmohammed/meetups/middlewares"
    "github.com/secmohammed/meetups/models"
)

func (a *mutationResolver) CreateAttendance(ctx context.Context, input models.CreateAttendanceInput) (*models.Attendee, error) {
    currentUser, _ := middlewares.GetCurrentUserFromContext(ctx)

    if _, err := a.MeetupsRepo.GetByID(input.MeetupID); err != nil {
        return nil, errors.New("meetup doens't exist")
    }
    attendee := &models.Attendee{
        Status:   input.Status.String(),
        MeetupID: input.MeetupID,
        UserID:   currentUser.ID,
    }
    return a.AttendeesRepo.Create(attendee)
}
func (a *mutationResolver) DeleteAttendance(ctx context.Context, id string) (bool, error) {
    currentUser, _ := middlewares.GetCurrentUserFromContext(ctx)

    attendee, err := a.AttendeesRepo.GetByID(id)
    if err != nil {
        return false, errors.New("Couldn't find this attendee to update")
    }
    if attendee.UserID != currentUser.ID {
        return false, errors.New("Unauthorized attempt")
    }
    return true, a.AttendeesRepo.Delete(attendee)

}
func (a *mutationResolver) UpdateAttendance(ctx context.Context, id string, status models.AttendanceStatus) (*models.Attendee, error) {
    currentUser, _ := middlewares.GetCurrentUserFromContext(ctx)

    attendee, err := a.AttendeesRepo.GetByID(id)
    if err != nil {
        return nil, errors.New("Couldn't find this attendee to update")
    }
    if attendee.UserID != currentUser.ID {
        return nil, errors.New("Unauthorized attempt")
    }

    attendee.Status = status.String()
    return a.AttendeesRepo.Update(attendee)
}
