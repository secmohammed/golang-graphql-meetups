//go:generate go run github.com/99designs/gqlgen -v

package resolvers

import (
    "encoding/json"
    "log"
    "sync"

    "github.com/go-pg/pg"
    "github.com/go-redis/redis"
    "github.com/nats-io/nats.go"
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
    mutex             sync.Mutex
    messageChannels   map[string]chan *models.Conversation
    userChannels      map[string]chan string
    redisClient       *utils.Cache
    nClient           *nats.EncodedConn
}

func NewResolver(DB *pg.DB, redisClient *utils.Cache) *Resolver {
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
        mutex:             sync.Mutex{},
        userChannels:      map[string]chan string{},
        messageChannels:   map[string]chan *models.Conversation{},
        redisClient:       redisClient,
        nClient:           nClient,
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
func (r *Resolver) createUser(user string) error {
    // Upsert user
    if err := r.redisClient.SAdd("users", user); err != nil {
        return err
    }
    // Notify new user joined
    r.mutex.Lock()
    for _, ch := range r.userChannels {
        ch <- user
    }
    r.mutex.Unlock()
    return nil
}
func (r *Resolver) StartSubscribingRedis() {
    log.Println("Start Subscribing Redis...")

    go func() {
        pubsub := r.redisClient.Subscribe("conversation")
        defer pubsub.Close()

        for {
            msgi, err := pubsub.Receive()
            if err != nil {
                panic(err)
            }

            switch msg := msgi.(type) {
            case *redis.Message:

                // Convert recieved string to Message.
                m := &models.Conversation{}
                if err := json.Unmarshal([]byte(msg.Payload), m); err != nil {
                    log.Println(err)
                    continue
                }

                // Notify new message.
                r.mutex.Lock()
                for _, ch := range r.messageChannels {
                    ch <- m
                }
                r.mutex.Unlock()

            default:
            }
        }
    }()
}
