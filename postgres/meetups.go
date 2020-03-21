package postgres

import (
    "time"

    "github.com/go-pg/pg"
    "github.com/go-pg/pg/orm"
    "github.com/secmohammed/meetups/models"
)

// MeetupsRepo is used to contain the db driver.
type MeetupsRepo struct {
    DB *pg.DB
}

//GetFilteredMeetupsBasedOnUser is used to fetch the meetups based on the filters such as
//start_date,end_date and name when provided
//in addition to the interests of the user.
func (m *MeetupsRepo) GetFilteredMeetupsBasedOnUser(userID string, filter *models.MeetupFilterInput, limit, offset *int) ([]*models.Meetup, error) {

    var meetups []*models.Meetup
    var categories []*models.Category
    // fetch the interests of user.
    err := m.DB.Model(&categories).Relation("Users", func(q *orm.Query) (*orm.Query, error) {
        return q.Where("category_user.user_id = ?", userID), nil
    }).Select()
    // pluck the ids of them.
    var ids []string
    for _, category := range categories {
        for _, user := range category.Users {
            if userID == user.ID {
                ids = append(ids, category.ID)
            }
        }
    }
    // select meetup which have these categories

    // filter query with specific id doesn't work due to that we are selecting all of the meetups first.
    query := m.DB.Model(&meetups).Relation("Categories", func(q *orm.Query) (*orm.Query, error) {
        return q.Where("category_meetup.category_id in (?)", pg.In(ids)), nil
    })
    if filter != nil && filter.Name != nil && *filter.Name != "" {
        query.Where("meetup.name = ?", filter.Name)
    }
    if filter != nil && filter.StartDate != nil && *filter.StartDate != "" {
        startDate, err := time.Parse("2006-01-02", *filter.StartDate)
        if err != nil {
            return nil, err
        }
        query.Where("start_date >= ?", startDate)
    }
    if filter != nil && filter.EndDate != nil && *filter.EndDate != "" {
        endDate, err := time.Parse("2006-01-02", *filter.EndDate)
        if err != nil {
            return nil, err
        }
        query.Where("end_date <= ?", endDate)
    }

    if limit != nil {
        query.Limit(*limit)
    }
    if offset != nil {
        query.Offset(*offset)
    }
    err = query.Select()
    if err != nil {
        return nil, err
    }
    // This is bullshit :)
    var result []*models.Meetup
    for _, meetup := range meetups {
        for _, category := range meetup.Categories {
            for _, id := range ids {
                if id == category.ID {
                    result = append(result, meetup)
                }
            }
        }
    }
    return result, nil

}

// GetMeetups is used to get meetups from database.
func (m *MeetupsRepo) GetMeetups(filter *models.MeetupFilterInput, limit, offset *int) ([]*models.Meetup, error) {
    var meetups []*models.Meetup
    query := m.DB.Model(&meetups).Order("id")
    if filter != nil && filter.Name != nil && *filter.Name != "" {
        query.Where("name = ?", filter.Name)

    }
    if filter != nil && filter.StartDate != nil && *filter.StartDate != "" {
        startDate, err := time.Parse("2006-01-02", *filter.StartDate)
        if err != nil {
            return nil, err
        }
        query.Where("start_date >= ?", startDate)
    }
    if filter != nil && filter.EndDate != nil && *filter.EndDate != "" {
        endDate, err := time.Parse("2006-01-02", *filter.EndDate)
        if err != nil {
            return nil, err
        }
        query.Where("end_date <= ?", endDate)
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
