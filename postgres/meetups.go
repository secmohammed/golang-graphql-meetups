package postgres

import (
    "github.com/go-pg/pg"
    "github.com/secmohammed/meetups/models"
)

// MeetupsRepo is used to contain the db driver.
type MeetupsRepo struct {
    DB *pg.DB
}

// GetMeetups is used to get meetups from database.
func (m *MeetupsRepo) GetMeetups() ([]*models.Meetup, error) {
    var meetups []*models.Meetup
    err := m.DB.Model(&meetups).Select()
    if err != nil {
        return nil, err
    }
    return meetups, nil
}

// CreateMeetup is used to create meetup for the database.
func (m *MeetupsRepo) CreateMeetup(meetup *models.Meetup) (*models.Meetup, error) {
    _, err := m.DB.Model(meetup).Returning("*").Insert()
    return meetup, err
}
