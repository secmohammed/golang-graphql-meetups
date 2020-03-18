package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/99designs/gqlgen/handler"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-pg/pg"
	_ "github.com/joho/godotenv/autoload"
	"github.com/rs/cors"
	"github.com/secmohammed/meetups/graphql"
	"github.com/secmohammed/meetups/graphql/loaders"
	"github.com/secmohammed/meetups/graphql/resolvers"
	"github.com/secmohammed/meetups/middlewares"
	"github.com/secmohammed/meetups/postgres"
	"github.com/secmohammed/meetups/utils"
)

func main() {
	DB := postgres.New(&pg.Options{
		User:     os.Getenv("DATABASE_USERNAME"),
		Password: os.Getenv("DATABASE_PASSWORD"),
		Database: os.Getenv("DATABASE_NAME"),
	})
	defer DB.Close()
	port := os.Getenv("APP_PORT")
	url := os.Getenv("APP_URL")
	userRepo := postgres.UsersRepo{DB: DB}
	router := chi.NewRouter()
	debug, _ := strconv.ParseBool(os.Getenv("APP_DEBUG"))
	cacheTTL, _ := strconv.ParseInt(os.Getenv("CACHE_TTL_IN_HOURS"), 10, 64)
	router.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{fmt.Sprintf("%s:%s", url, port)},
		Debug:            debug,
		AllowCredentials: true,
	}).Handler)
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middlewares.AuthMiddleware(userRepo))
	c := graphql.Config{Resolvers: &resolvers.Resolver{
		MeetupsRepo: postgres.MeetupsRepo{DB: DB},
		UsersRepo:   userRepo,
	}}
	cache, err := utils.NewCache(os.Getenv("REDIS_ADDRESS"), time.Duration(cacheTTL)*time.Hour)
	if err != nil {
		log.Fatalf("Cannot create APQ redis cache: %v", err)
	}
	queryHandler := handler.GraphQL(graphql.NewExecutableSchema(c), handler.EnablePersistedQueryCache(cache))

	router.Handle("/", handler.Playground("GraphQL playground", "/query"))
	router.Handle("/query", loaders.DataloaderMiddleware(DB, queryHandler))

	log.Printf("connect to %s:%s/ for GraphQL playground", url, port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
