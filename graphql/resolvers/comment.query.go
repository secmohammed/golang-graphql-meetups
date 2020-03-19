package resolvers

import (
    "context"

    "github.com/secmohammed/meetups/graphql"
    "github.com/secmohammed/meetups/models"
)

func (r *Resolver) Comment() graphql.CommentResolver {
    return &commentResolver{r}
}
func (c *queryResolver) Comments(ctx context.Context, meetupID string) ([]*models.Comment, error) {
    return c.CommentsRepo.GetCommentsForMeetup(meetupID)
}
