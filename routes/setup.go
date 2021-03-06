package routes

import (
    "context"
    "fmt"
    "log"
    "net/http"
    "os"
    "strconv"
    "time"

    packageGraphQL "github.com/99designs/gqlgen/graphql"
    "github.com/99designs/gqlgen/handler"
    "github.com/go-chi/chi"
    "github.com/go-chi/chi/middleware"
    "github.com/go-pg/pg"
    "github.com/gorilla/websocket"
    "github.com/prprprus/scheduler"
    "github.com/rs/cors"
    "github.com/secmohammed/meetups/graphql"
    "github.com/secmohammed/meetups/graphql/loaders"
    "github.com/secmohammed/meetups/graphql/resolvers"
    "github.com/secmohammed/meetups/middlewares"
    "github.com/secmohammed/meetups/models"
    "github.com/secmohammed/meetups/postgres"
    "github.com/secmohammed/meetups/utils"
    "github.com/secmohammed/meetups/utils/errors"
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
    cache, err := utils.NewCache(os.Getenv("REDIS_ADDRESS"), time.Duration(cacheTTL)*time.Hour)
    if err != nil {
        log.Fatalf("Cannot create APQ redis cache: %v", err)
    }
    s, err := scheduler.NewScheduler(1000)
    if err != nil {
        log.Fatalf("Cannot create scheduler: %v", err)
    }
    resolver := resolvers.NewResolver(DB, cache, s)
    c := graphql.Config{Resolvers: resolver}
    c.Directives.Authentication = func(ctx context.Context, obj interface{}, next packageGraphQL.Resolver, auth models.Authentication) (res interface{}, err error) {
        _, err = middlewares.GetCurrentUserFromContext(ctx)
        if err != nil && auth == "AUTHENTICATED" {
            return nil, errors.ErrUnauthenticated
        }
        if err == nil && auth == "GUEST" {
            return nil, errors.ErrAuthenticated
        }
        return next(ctx)
    }
    c.Directives.Can = func(ctx context.Context, obj interface{}, next packageGraphQL.Resolver, role string) (res interface{}, err error) {
        _, err = middlewares.GetCurrentUserFromContext(ctx)
        ok, err := middlewares.HasRole(ctx, role, userRepo)
        if err != nil || !ok {
            return nil, errors.ErrInsufficientPermissions
        }

        return next(ctx)
    }
    router.Use(middlewares.AuthMiddleware(userRepo))

    queryHandler := handler.GraphQL(
        graphql.NewExecutableSchema(c),
        handler.EnablePersistedQueryCache(cache),
        handler.ComplexityLimit(200),
        handler.WebsocketUpgrader(
            websocket.Upgrader{
                CheckOrigin: func(r *http.Request) bool {
                    return true
                },
            }),
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
