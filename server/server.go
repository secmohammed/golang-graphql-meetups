package main

import (
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/handler"
	"github.com/go-pg/pg"
	"github.com/secmohammed/meetups/graphql"
	"github.com/secmohammed/meetups/graphql/loaders"
	"github.com/secmohammed/meetups/graphql/resolvers"
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

	c := graphql.Config{Resolvers: &resolvers.Resolver{
		MeetupsRepo: postgres.MeetupsRepo{DB: DB},
		UsersRepo:   postgres.UsersRepo{DB: DB},
	}}
	queryHandler := handler.GraphQL(graphql.NewExecutableSchema(c))

	http.Handle("/", handler.Playground("GraphQL playground", "/query"))
	http.Handle("/query", loaders.DataloaderMiddleware(DB, queryHandler))

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
