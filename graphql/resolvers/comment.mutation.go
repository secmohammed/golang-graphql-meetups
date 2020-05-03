package resolvers

import (
    "context"

    "github.com/secmohammed/meetups/middlewares"
    "github.com/secmohammed/meetups/models"
    "github.com/secmohammed/meetups/utils/errors"
)

type commentResolver struct{ *Resolver }

func (c *mutationResolver) CreateComment(ctx context.Context, input models.CreateCommentInput) (*models.Comment, error) {
    currentUser, _ := middlewares.GetCurrentUserFromContext(ctx)

    if _, err := c.MeetupsRepo.GetByID(input.MeetupID); err != nil {
        return nil, errors.ErrRecordNotFound
    }

    if err := input.Validate(); err != nil {
        return nil, err
    }

    comment := &models.Comment{
        Body:     input.Body,
        UserID:   currentUser.ID,
        MeetupID: input.MeetupID,
    }
    if input.ParentID != "" {
        //TODO: if comment has a parent_id, fetch the user_id that's attached to this record and notify this user.
        foundComment, err := c.CommentsRepo.GetByID(input.ParentID)
        if err != nil {
            return nil, errors.ErrRecordNotFound
        }
        // create notification for this user.
        notification := &models.Notification{
            UserID:         foundComment.UserID,
            NotifiableType: "comment_created",
            NotifiableID:   foundComment.ID,
        }
        notification, err = c.NotificationsRepo.Create(notification)
        if err != nil {
            return nil, errors.ErrInternalError
        }
        c.nClient.Publish("notification.user_"+foundComment.UserID, notification)
        comment.ParentID = input.ParentID
    }
    if input.GroupID != "" {
        _, err := c.GroupsRepo.GetByID(input.GroupID)
        if err != nil {
            return nil, errors.ErrRecordNotFound
        }

        comment.GroupID = input.GroupID
    }

    return c.CommentsRepo.Create(comment)

}
func (c *mutationResolver) UpdateComment(ctx context.Context, id string, input models.UpdateCommentInput) (*models.Comment, error) {
    currentUser, _ := middlewares.GetCurrentUserFromContext(ctx)

    comment, err := c.CommentsRepo.GetByID(id)
    if err != nil {
        return nil, errors.ErrRecordNotFound
    }
    if comment.UserID != currentUser.ID {
        return nil, errors.ErrUnauthenticated
    }
    if err := input.Validate(); err != nil {
        return nil, err
    }

    comment.Body = input.Body

    return c.CommentsRepo.Update(comment)

}
func (c *mutationResolver) DeleteComment(ctx context.Context, id string) (bool, error) {
    currentUser, _ := middlewares.GetCurrentUserFromContext(ctx)

    comment, err := c.CommentsRepo.GetByID(id)
    if err != nil {
        return false, errors.ErrRecordNotFound
    }
    if comment.UserID != currentUser.ID {
        return false, errors.ErrUnauthenticated
    }
    return true, c.CommentsRepo.Delete(comment)
}
