//go:generate go run github.com/99designs/gqlgen -v

package meetmeup

import (
	"context"
	"errors"

	"github.com/secmohammed/meetups/models"
	"github.com/secmohammed/meetups/postgres"
)

type Resolver struct {
	MeetupsRepo postgres.MeetupsRepo
	UsersRepo   postgres.UsersRepo
}

type queryResolver struct{ *Resolver }

func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}

type mutationResolver struct{ *Resolver }

func (r *Resolver) Mutation() MutationResolver {
	return &mutationResolver{r}
}

type meetupResolver struct{ *Resolver }

func (m *meetupResolver) User(ctx context.Context, obj *models.Meetup) (*models.User, error) {

	return m.Resolver.UsersRepo.GetByID(obj.UserID)
}

func (r *Resolver) Meetup() MeetupResolver {
	return &meetupResolver{r}
}

type userResolver struct{ *Resolver }

func (r *Resolver) User() UserResolver {
	return &userResolver{r}
}

func (u *userResolver) Meetups(ctx context.Context, obj *models.User) ([]*models.Meetup, error) {
	return nil, nil
}

func (m *mutationResolver) CreateMeetup(ctx context.Context, input NewMeetup) (*models.Meetup, error) {
	if len(input.Name) < 3 {
		return nil, errors.New("Name must be more than 3 characters")
	}
	if len(input.Description) < 3 {
		return nil, errors.New("Description must be more than 3 characters")
	}
	meetup := &models.Meetup{
		Name:        input.Name,
		Description: input.Description,
		UserID:      "1",
	}
	return m.MeetupsRepo.CreateMeetup(meetup)
}

func (r *queryResolver) Meetups(ctx context.Context) ([]*models.Meetup, error) {

	return r.MeetupsRepo.GetMeetups()
}
