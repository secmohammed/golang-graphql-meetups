package routes

import (
    "fmt"
    "log"
    "os"
    "strconv"
    "time"

    "github.com/99designs/gqlgen/handler"
    "github.com/go-chi/chi"
    "github.com/go-chi/chi/middleware"
    "github.com/go-pg/pg"
    "github.com/rs/cors"
    "github.com/secmohammed/meetups/graphql"
    "github.com/secmohammed/meetups/graphql/loaders"
    "github.com/secmohammed/meetups/graphql/resolvers"
    "github.com/secmohammed/meetups/middlewares"
    "github.com/secmohammed/meetups/postgres"
    "github.com/secmohammed/meetups/utils"
)

var (
    debug, _    = strconv.ParseBool(os.Getenv("APP_DEBUG"))
    cacheTTL, _ = strconv.ParseInt(os.Getenv("CACHE_TTL_IN_HOURS"), 10, 64)
    port        = os.Getenv("APP_PORT")
    url         = os.Getenv("APP_URL")
)

//SetupRoutes is used to setup the routes for the server.
func SetupRoutes(DB *pg.DB) *chi.Mux {
    userRepo := postgres.UsersRepo{DB: DB}

    router := chi.NewRouter()
    c := graphql.Config{Resolvers: &resolvers.Resolver{
        MeetupsRepo:    postgres.MeetupsRepo{DB: DB},
        UsersRepo:      userRepo,
        CommentsRepo:   postgres.CommentsRepo{DB: DB},
        CategoriesRepo: postgres.CategoriesRepo{DB: DB},
        AttendeesRepo:  postgres.AttendeesRepo{DB: DB},
    }}
    router.Use(middlewares.AuthMiddleware(userRepo))

    cache, err := utils.NewCache(os.Getenv("REDIS_ADDRESS"), time.Duration(cacheTTL)*time.Hour)
    if err != nil {
        log.Fatalf("Cannot create APQ redis cache: %v", err)
    }
    queryHandler := handler.GraphQL(
        graphql.NewExecutableSchema(c),
        handler.EnablePersistedQueryCache(cache),
        handler.ComplexityLimit(10),
    )

    router.Use(cors.New(cors.Options{
        AllowedOrigins:   []string{fmt.Sprintf("%s:%s", url, port)},
        Debug:            debug,
        AllowCredentials: true,
    }).Handler)
    router.Use(middleware.RequestID)
    router.Use(middleware.Logger)
    router.Handle("/", handler.Playground("GraphQL playground", "/query"))
    router.Handle("/query", loaders.DataloaderMiddleware(DB, queryHandler))
    return router

}
