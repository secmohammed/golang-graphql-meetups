package resolvers

import (
    "context"
    "errors"

    "github.com/secmohammed/meetups/middlewares"
    "github.com/secmohammed/meetups/models"
)

type commentResolver struct{ *Resolver }

func (c *mutationResolver) CreateComment(ctx context.Context, input models.CreateCommentInput) (*models.Comment, error) {
    currentUser, err := middlewares.GetCurrentUserFromContext(ctx)
    if err != nil {
        return nil, ErrUnauthenticated
    }

    if _, err := c.MeetupsRepo.GetByID(input.MeetupID); err != nil {
        return nil, errors.New("meetup doens't exist")
    }

    if err := input.Validate(); err != nil {
        return nil, err
    }
    comment := &models.Comment{
        Body:     input.Body,
        UserID:   currentUser.ID,
        MeetupID: input.MeetupID,
        ParentID: input.ParentID,
    }

    return c.CommentsRepo.Create(comment)

}
func (c *mutationResolver) UpdateComment(ctx context.Context, id string, input models.UpdateCommentInput) (*models.Comment, error) {
    currentUser, err := middlewares.GetCurrentUserFromContext(ctx)
    if err != nil {
        return nil, ErrUnauthenticated
    }
    comment, err := c.CommentsRepo.GetByID(id)
    if err != nil {
        return nil, errors.New("Couldn't find this comment")
    }
    if comment.UserID != currentUser.ID {
        return nil, errors.New("Unauthorized Attempt")
    }
    if err := input.Validate(); err != nil {
        return nil, err
    }
    comment.Body = input.Body
    comment.ParentID = input.ParentID

    return c.CommentsRepo.Update(comment)

}
func (c *mutationResolver) DeleteComment(ctx context.Context, id string) (bool, error) {
    currentUser, err := middlewares.GetCurrentUserFromContext(ctx)
    if err != nil {
        return false, ErrUnauthenticated
    }
    comment, err := c.CommentsRepo.GetByID(id)
    if err != nil {
        return false, errors.New("Couldn't find this comment")
    }
    if comment.UserID != currentUser.ID {
        return false, errors.New("Unauthorized Attempt")
    }
    return true, c.CommentsRepo.Delete(comment)
}
