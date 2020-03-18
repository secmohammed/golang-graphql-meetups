package postgres

import (
    "github.com/go-pg/pg"
    "github.com/secmohammed/meetups/models"
)

// CommentsRepo is used to contain the db driver.
type CommentsRepo struct {
    DB *pg.DB
}

// GetCommentsForUser is used to get comments for the passed user by its id.
func (c *CommentsRepo) GetCommentsForUser(user *models.User) ([]*models.Comment, error) {
    var comments []*models.Comment
    err := c.DB.Model(&comments).Where("user_id = ? ", user.ID).Order("id").Select()
    return comments, err
}

// GetCommentsForMeetup is used to get meetups for the passed user by its id.
func (c *CommentsRepo) GetCommentsForMeetup(meetup *models.Meetup) ([]*models.Comment, error) {
    var comments []*models.Comment
    err := c.DB.Model(&comments).Where("meetup_id = ? ", meetup.ID).Order("id").Select()
    if err != nil {
        return nil, err
    }

    return comments, err
}
