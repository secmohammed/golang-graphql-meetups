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
func (m *MeetupsRepo) GetMeetups(filter *models.MeetupFilter, limit, offset *int) ([]*models.Meetup, error) {
    var meetups []*models.Meetup
    query := m.DB.Model(&meetups).Order("id")
    if filter != nil && filter.Name != nil && *filter.Name != "" {

    }
    if limit != nil {
        query.Limit(*limit)
    }
    if offset != nil {
        query.Offset(*offset)
    }
    err := query.Select()
    if err != nil {
        return nil, err
    }
    return meetups, nil
}

// Create is used to create meetup for the database.
func (m *MeetupsRepo) Create(meetup *models.Meetup) (*models.Meetup, error) {
    _, err := m.DB.Model(meetup).Returning("*").Insert()
    return meetup, err
}

// Update is used to update the passed meetup by id.
func (m *MeetupsRepo) Update(meetup *models.Meetup) (*models.Meetup, error) {
    _, err := m.DB.Model(meetup).Where("id = ?", meetup.ID).Update()
    return meetup, err
}

// GetByID is used to fetch meetup by id.
func (m *MeetupsRepo) GetByID(id string) (*models.Meetup, error) {
    meetup := models.Meetup{}
    err := m.DB.Model(&meetup).Where("id = ?", id).Select()
    if err != nil {
        return nil, err
    }
    return &meetup, nil
}

// Delete is used to delete meetup by its id.
func (m *MeetupsRepo) Delete(meetup *models.Meetup) error {
    _, err := m.DB.Model(meetup).Where("id = ?", meetup.ID).Delete()
    return err
}

// GetMeetupsForUser is used to get meetups for the passed user by its id.
func (m *MeetupsRepo) GetMeetupsForUser(user *models.User) ([]*models.Meetup, error) {
    var meetups []*models.Meetup
    err := m.DB.Model(&meetups).Where("user_id = ? ", user.ID).Order("id").Select()
    return meetups, err
}
