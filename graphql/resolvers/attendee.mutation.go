package resolvers

import (
    "context"
    "time"

    "github.com/prprprus/scheduler"
    "github.com/secmohammed/meetups/mails"
    "github.com/secmohammed/meetups/middlewares"
    "github.com/secmohammed/meetups/models"
    "github.com/secmohammed/meetups/utils/errors"
)

func (a *mutationResolver) CreateAttendance(ctx context.Context, input models.CreateAttendanceInput) (*models.Attendee, error) {
    currentUser, _ := middlewares.GetCurrentUserFromContext(ctx)
    meetup, err := a.MeetupsRepo.GetByID(input.MeetupID)
    if err != nil {
        return nil, errors.ErrRecordNotFound
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
    // TODO: to Delete attendance, it must be like sort of cancelation because a user mustn't delete his attendance of a meetup he attended.
    attendee, err := a.AttendeesRepo.GetByID(id)
    if err != nil {
        return false, errors.ErrRecordNotFound
    }
    if attendee.UserID != currentUser.ID {
        return false, errors.ErrUnauthenticated
    }
    // TODO: If the meetup hasn't been made yet, we must delete the reminder for the user.
    // to do so, we must store the schedule id that's made at redis.
    return true, a.AttendeesRepo.Delete(attendee)

}
func (a *mutationResolver) UpdateAttendance(ctx context.Context, id string, status models.AttendanceStatus) (*models.Attendee, error) {
    currentUser, _ := middlewares.GetCurrentUserFromContext(ctx)
    // TODO: user can't update a status if the meeeting has passed away.
    attendee, err := a.AttendeesRepo.GetByID(id)
    if err != nil {
        return nil, errors.ErrRecordNotFound
    }
    if attendee.UserID != currentUser.ID {
        return nil, errors.ErrUnauthenticated
    }

    attendee.Status = status.String()
    return a.AttendeesRepo.Update(attendee)
}
