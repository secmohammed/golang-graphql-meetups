package main

import (
	"log"
	"net/http"
	"os"

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
)

const defaultPort = "8080"

func main() {
	DB := postgres.New(&pg.Options{
		User:     "mohammed",
		Password: "root",
		Database: "meetups",
	})
	defer DB.Close()
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}
	userRepo := postgres.UsersRepo{DB: DB}
	router := chi.NewRouter()
	router.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:8080"},
		Debug:            true,
		AllowCredentials: true,
	}).Handler)
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middlewares.AuthMiddleware(userRepo))
	c := graphql.Config{Resolvers: &resolvers.Resolver{
		MeetupsRepo: postgres.MeetupsRepo{DB: DB},
		UsersRepo:   userRepo,
	}}
	queryHandler := handler.GraphQL(graphql.NewExecutableSchema(c))

	router.Handle("/", handler.Playground("GraphQL playground", "/query"))
	router.Handle("/query", loaders.DataloaderMiddleware(DB, queryHandler))

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
