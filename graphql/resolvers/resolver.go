//go:generate go run github.com/99designs/gqlgen -v

package resolvers

import (
	"github.com/secmohammed/meetups/postgres"
)

type Resolver struct {
	MeetupsRepo postgres.MeetupsRepo
	UsersRepo   postgres.UsersRepo
}
