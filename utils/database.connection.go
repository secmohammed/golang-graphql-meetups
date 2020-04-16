package utils

import (
    "os"
    "sync"

    "github.com/go-pg/pg"
    "github.com/secmohammed/meetups/postgres"
)

var once sync.Once
var db *pg.DB

func NewDatabaseConnection() *pg.DB {
    once.Do(func() {
        db = postgres.New(&pg.Options{
            User:     os.Getenv("DATABASE_USERNAME"),
            Password: os.Getenv("DATABASE_PASSWORD"),
            Database: os.Getenv("DATABASE_NAME"),
        })
    })
    return db
}
