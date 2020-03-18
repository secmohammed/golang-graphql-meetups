package utils

import (
    "os"

    "github.com/go-pg/pg"
    "github.com/secmohammed/meetups/postgres"
)

func SetDatabaseConnection() *pg.DB {
    return postgres.New(&pg.Options{
        User:     os.Getenv("DATABASE_USERNAME"),
        Password: os.Getenv("DATABASE_PASSWORD"),
        Database: os.Getenv("DATABASE_NAME"),
    })
}
