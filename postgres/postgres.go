package postgres

import (
    "fmt"

    "github.com/go-pg/pg"
)

type dbLogger struct{}

func (d dbLogger) BeforeQuery(q *pg.QueryEvent) {
}

func (d dbLogger) AfterQuery(q *pg.QueryEvent) {
    fmt.Println(q.FormattedQuery())
}

// New function is used to create a new connection of pgsql
func New(options *pg.Options) *pg.DB {

    db := pg.Connect(options)
    db.AddQueryHook(dbLogger{})
    return db
}
