package postgres

import (
    "github.com/go-pg/pg"
    "github.com/secmohammed/meetups/models"
)

// AttendeesRepo is used to contain the db driver.
type AttendeesRepo struct {
    DB *pg.DB
}

//CreateAttendanceForUser is used to create attendance for user on specific meetup.
func (a *AttendeesRepo) CreateAttendanceForUser() (*models.Attendee, error) {
    return nil, nil
}
