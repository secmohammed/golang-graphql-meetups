package postgres

import (
    "github.com/go-pg/pg"
    "github.com/secmohammed/meetups/models"
)

// CommentsRepo is used to contain the db driver.
type CommentsRepo struct {
    DB *pg.DB
}

//Create is used to create a comment using the passed struct.
func (c *CommentsRepo) Create(comment *models.Comment) (*models.Comment, error) {
    _, err := c.DB.Model(comment).Returning("*").Insert()
    return comment, err
}

// Update is used to update the passed meetup by id.
func (c *CommentsRepo) Update(comment *models.Comment) (*models.Comment, error) {
    _, err := c.DB.Model(comment).Where("id = ?", comment.ID).Update()
    return comment, err
}

// GetByID is used to fetch meetup by id.
func (c *CommentsRepo) GetByID(id string) (*models.Comment, error) {
    comment := models.Comment{}
    err := c.DB.Model(&comment).Where("id = ?", id).Select()
    if err != nil {
        return nil, err
    }
    return &comment, nil
}

// Delete is used to delete meetup by its id.
func (c *CommentsRepo) Delete(comment *models.Comment) error {
    _, err := c.DB.Model(comment).Where("id = ?", comment.ID).Delete()
    return err
}

// GetCommentsForUser is used to get comments for the passed user by its id.
func (c *CommentsRepo) GetCommentsForUser(user *models.User) ([]*models.Comment, error) {
    var comments []*models.Comment
    err := c.DB.Model(&comments).Where("user_id = ? ", user.ID).Where("parent_id = ?", nil).Order("id").Select()
    return comments, err
}

//GetRepliesForComment is used to fetch replies for the passed comment id.
func (c *CommentsRepo) GetRepliesForComment(id string) ([]*models.Comment, error) {
    var comments []*models.Comment
    err := c.DB.Model(&comments).Where("parent_id = ?", id).Order("id").Select()
    return comments, err
}

// GetCommentsForMeetup is used to get meetups for the passed user by its id.
func (c *CommentsRepo) GetCommentsForMeetup(id string) ([]*models.Comment, error) {
    var comments []*models.Comment
    err := c.DB.Model(&comments).Where("meetup_id = ?", id).Where("parent_id IS NULL").Order("id").Select()
    if err != nil {
        return nil, err
    }

    return comments, err
}
