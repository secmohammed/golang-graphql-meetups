package resolvers

import (
    "context"
    "fmt"
    "log"
    "os"

    "github.com/appleboy/go-fcm"
    "github.com/secmohammed/meetups/middlewares"
    "github.com/secmohammed/meetups/models"
    "github.com/secmohammed/meetups/utils/errors"
)

type commentResolver struct{ *Resolver }

func (c *mutationResolver) CreateComment(ctx context.Context, input models.CreateCommentInput) (*models.Comment, error) {
    currentUser, _ := middlewares.GetCurrentUserFromContext(ctx)
    meetup, err := c.MeetupsRepo.GetByID(input.MeetupID)
    if err != nil {
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
        //if comment has a parent_id, fetch the user_id that's attached to this record and notify this user.
        foundComment, err := c.CommentsRepo.GetByID(input.ParentID)
        if err != nil {
            return nil, errors.ErrRecordNotFound
        }
        // create notification for this user.
        notification := &models.Notification{
            UserID:         foundComment.UserID, // the one we wish to notify
            NotifiableType: "reply_created",
            NotifiableID:   foundComment.ID, // the concerened comment that caused the notification
        }
        notification, err = c.NotificationsRepo.Create(notification)
        if err != nil {
            return nil, errors.ErrInternalError
        }
        c.nClient.Publish("notification.user_"+foundComment.UserID, notification)
        comment.ParentID = input.ParentID
        msg := &fcm.Message{
            To: foundComment.User.DeviceToken,
            Data: map[string]interface{}{
                "comment": comment,
                "message": fmt.Sprintf("%s replied  you with: %s", foundComment.User.Username, comment.Body),
            },
        }
        client, err := fcm.NewClient(os.Getenv("fcm_token"))
        if err != nil {
            log.Fatalln(err)
        }

        // Send the message and receive the response without retries.
        _, err = client.Send(msg)
        if err != nil {
            log.Fatalln(err)
        }
    }
    if input.GroupID != "" {
        _, err := c.GroupsRepo.GetByID(input.GroupID)
        if err != nil {
            return nil, errors.ErrRecordNotFound
        }

        comment.GroupID = input.GroupID
    }
    notification := &models.Notification{
        UserID:         meetup.UserID, // the one we wish to notify which is the meetup owner.
        NotifiableType: "comment_created",
        NotifiableID:   meetup.ID, // the concerened meetup that caused the notification.
    }
    notification, err = c.NotificationsRepo.Create(notification)
    if err != nil {
        return nil, errors.ErrInternalError
    }
    c.nClient.Publish("notification.user_"+meetup.UserID, notification)

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
