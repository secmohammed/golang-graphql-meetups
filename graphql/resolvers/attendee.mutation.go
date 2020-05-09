package resolvers

import (
    "context"
    "time"

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
    // disallow users to create attendance for meetups that's already after meetup end date.
    if time.Now().After(meetup.EndDate) {
        return nil, errors.ErrUnauthenticated
    }

    attendee := &models.Attendee{
        Status:   input.Status.String(),
        MeetupID: input.MeetupID,
        UserID:   currentUser.ID,
    }
    delayedTimeInSeconds := int(meetup.StartDate.Add(time.Duration(-1) * time.Hour).Unix())
    jobID := a.scheduler.Delay().Second(delayedTimeInSeconds).Do(mails.SendReminderEmailToAttendee, currentUser, meetup)
    a.redisClient.Add(ctx, "job_id_for_meetup_"+meetup.ID, jobID)
    return a.AttendeesRepo.Create(attendee)
}
func (a *mutationResolver) DeleteAttendance(ctx context.Context, id string) (bool, error) {
    currentUser, _ := middlewares.GetCurrentUserFromContext(ctx)
    attendee, err := a.AttendeesRepo.GetByID(id, "Meetup")
    if err != nil {
        return false, errors.ErrRecordNotFound
    }
    if attendee.UserID != currentUser.ID {
        return false, errors.ErrUnauthenticated
    }
    // If the meetup hasn't been made yet, we must delete the reminder for the user.
    // to do so, we must store the schedule id that's made at redis.
    if attendee.CreatedAt.After(attendee.Meetup.StartDate) {
        return false, errors.ErrUnauthenticated
    }
    jobID, ok := a.redisClient.Get(ctx, "job_id_for_meetup_"+attendee.MeetupID)
    if !ok {
        return false, errors.ErrRecordNotFound
    }
    a.scheduler.CancelJob(jobID)
    return true, a.AttendeesRepo.Delete(attendee)

}
func (a *mutationResolver) UpdateAttendance(ctx context.Context, id string, status models.AttendanceStatus) (*models.Attendee, error) {
    currentUser, _ := middlewares.GetCurrentUserFromContext(ctx)
    // user can't update a status if the meeeting has passed away.
    attendee, err := a.AttendeesRepo.GetByID(id, "Meetup")
    if time.Now().After(attendee.Meetup.StartDate) {
        return nil, errors.ErrUnauthenticated
    }

    if err != nil {
        return nil, errors.ErrRecordNotFound
    }
    if attendee.UserID != currentUser.ID {
        return nil, errors.ErrUnauthenticated
    }

    attendee.Status = status.String()
    return a.AttendeesRepo.Update(attendee)
}
