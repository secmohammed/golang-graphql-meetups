package postgres

import (
    "github.com/go-pg/pg"
    "github.com/secmohammed/meetups/models"
)

// AttendeesRepo is used to contain the db driver.
type AttendeesRepo struct {
    DB *pg.DB
}

//Create is used to create a attendee using the passed struct.
func (a *AttendeesRepo) Create(attendee *models.Attendee) (*models.Attendee, error) {
    _, err := a.DB.Model(attendee).Returning("*").Insert()
    return attendee, err
}

// GetByID is used to fetch meetup by id.
func (a *AttendeesRepo) GetByID(id, relation string) (*models.Attendee, error) {
    attendee := models.Attendee{}
    query := a.DB.Model(&attendee).Where("id = ?", id)
    if relation != "" {
        query.Relation(relation)
    }
    err := query.First()
    if err != nil {
        return nil, err
    }
    return &attendee, nil
}

// Update is used to update the passed attendee by id.
func (a *AttendeesRepo) Update(attendee *models.Attendee) (*models.Attendee, error) {
    _, err := a.DB.Model(attendee).Where("id = ?", attendee.ID).Update()
    return attendee, err
}

// Delete is used to delete attendee by its id.
func (a *AttendeesRepo) Delete(attendee *models.Attendee) error {
    _, err := a.DB.Model(attendee).Where("id = ?", attendee.ID).Delete()
    return err
}
