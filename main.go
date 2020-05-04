package main

import (
	"log"
	"net/http"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"github.com/secmohammed/meetups/routes"
	"github.com/secmohammed/meetups/utils"
)

var (
	port = os.Getenv("APP_PORT")
	url  = os.Getenv("APP_URL")
)

func main() {
	DB := utils.NewDatabaseConnection()
	defer DB.Close()
	router := routes.SetupRoutes(DB)

	log.Printf("connect to %s:%s/ for GraphQL playground", url, port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
