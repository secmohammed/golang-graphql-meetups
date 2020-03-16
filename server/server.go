package main

import (
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/handler"
	"github.com/go-pg/pg"
	meetmeup "github.com/secmohammed/meetups"
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

	c := meetmeup.Config{Resolvers: &meetmeup.Resolver{
		MeetupsRepo: postgres.MeetupsRepo{DB: DB},
		UsersRepo:   postgres.UsersRepo{DB: DB},
	}}

	http.Handle("/", handler.Playground("GraphQL playground", "/query"))
	http.Handle("/query", handler.GraphQL(meetmeup.NewExecutableSchema(c)))

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
