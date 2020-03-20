package resolvers

import (
    "context"

    "github.com/secmohammed/meetups/graphql"
    "github.com/secmohammed/meetups/graphql/loaders"
    "github.com/secmohammed/meetups/models"
)

func (r *Resolver) Comment() graphql.CommentResolver {
    return &commentResolver{r}
}
func (c *queryResolver) Comments(ctx context.Context, meetupID string) ([]*models.Comment, error) {
    return c.CommentsRepo.GetCommentsForMeetup(meetupID)
}
func (c *commentResolver) User(ctx context.Context, obj *models.Comment) (*models.User, error) {
    return loaders.GetLoaders(ctx).UserByID.Load(obj.UserID)

}
func (c *commentResolver) Replies(ctx context.Context, obj *models.Comment) ([]*models.Comment, error) {
    return c.CommentsRepo.GetRepliesForComment(obj.ID)
}
