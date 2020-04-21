package resolvers

import (
    "context"
    "errors"
    "time"

    "github.com/prprprus/scheduler"
    "github.com/secmohammed/meetups/mails"
    "github.com/secmohammed/meetups/middlewares"
    "github.com/secmohammed/meetups/models"
)

func (a *mutationResolver) CreateAttendance(ctx context.Context, input models.CreateAttendanceInput) (*models.Attendee, error) {
    currentUser, _ := middlewares.GetCurrentUserFromContext(ctx)
    meetup, err := a.MeetupsRepo.GetByID(input.MeetupID)
    if err != nil {
        return nil, errors.New("meetup doens't exist")
    }
    attendee := &models.Attendee{
        Status:   input.Status.String(),
        MeetupID: input.MeetupID,
        UserID:   currentUser.ID,
    }
    s, err := scheduler.NewScheduler(1000)
    parsed, err := time.Parse("2006-01-02 15:04:05+02", meetup.StartDate)
    if err != nil {
        return nil, err
    }
    delayedTimeInSeconds := int(parsed.Add(time.Duration(-1) * time.Hour).Unix())
    s.Delay().Second(delayedTimeInSeconds).Do(mails.SendReminderEmailToAttendee, currentUser, meetup)
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
    // TODO: If the meetup hasn't been made yet, we must delete the reminder for the user.
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
