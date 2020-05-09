//go:generate go run github.com/99designs/gqlgen -v

package resolvers

import (
    "log"

    "github.com/go-pg/pg"
    "github.com/nats-io/nats.go"
    "github.com/prprprus/scheduler"
    "github.com/secmohammed/meetups/graphql"
    "github.com/secmohammed/meetups/models"
    "github.com/secmohammed/meetups/postgres"
    "github.com/secmohammed/meetups/utils"
)

//Resolver struct.
type Resolver struct {
    MeetupsRepo       postgres.MeetupsRepo
    UsersRepo         postgres.UsersRepo
    CommentsRepo      postgres.CommentsRepo
    CategoriesRepo    postgres.CategoriesRepo
    AttendeesRepo     postgres.AttendeesRepo
    ConversationsRepo postgres.ConversationsRepo
    GroupsRepo        postgres.GroupsRepo
    RolesRepo         postgres.RolesRepo
    NotificationsRepo postgres.NotificationsRepo
    messageChannels   map[string]chan *models.Conversation
    nClient           *nats.EncodedConn
    redisClient       *utils.Cache
    scheduler         *scheduler.Scheduler
}

func NewResolver(DB *pg.DB, redisClient *utils.Cache, s *scheduler.Scheduler) *Resolver {
    nc, err := nats.Connect(nats.DefaultURL)
    nClient, _ := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
    if err != nil {
        log.Fatalln(err)
    }
    return &Resolver{
        MeetupsRepo:       postgres.MeetupsRepo{DB: DB},
        UsersRepo:         postgres.UsersRepo{DB: DB},
        CommentsRepo:      postgres.CommentsRepo{DB: DB},
        CategoriesRepo:    postgres.CategoriesRepo{DB: DB},
        AttendeesRepo:     postgres.AttendeesRepo{DB: DB},
        ConversationsRepo: postgres.ConversationsRepo{DB: DB},
        GroupsRepo:        postgres.GroupsRepo{DB: DB},
        NotificationsRepo: postgres.NotificationsRepo{DB: DB},
        RolesRepo:         postgres.RolesRepo{DB: DB},
        messageChannels:   map[string]chan *models.Conversation{},
        nClient:           nClient,
        redisClient:       redisClient,
        scheduler:         s,
    }
}

type mutationResolver struct{ *Resolver }

//
type queryResolver struct{ *Resolver }

// Mutation method is used to resolve the mutations
func (r *Resolver) Mutation() graphql.MutationResolver {
    return &mutationResolver{r}
}

// Query method is used to resolve the queries
func (r *Resolver) Query() graphql.QueryResolver {
    return &queryResolver{r}
}
