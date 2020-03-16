package loaders

import (
    "context"
    "net/http"
    "time"

    "github.com/go-pg/pg"
    "github.com/secmohammed/meetups/models"
    "github.com/secmohammed/meetups/utils"
)

const userloaderKey = utils.ContextKey("userloader")

//DataloaderMiddleware is used to load the data related to users.
func DataloaderMiddleware(db *pg.DB, next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        userloader := UserLoader{
            maxBatch: 100,
            wait:     1 * time.Millisecond,
            fetch: func(ids []string) ([]*models.User, []error) {
                var users []*models.User
                err := db.Model(&users).Where("id in (?)", pg.In(ids)).Select()
                return users, []error{err}
            },
        }
        ctx := context.WithValue(r.Context(), userloaderKey, &userloader)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

//GetUserLoader is used to fetch the user loader instance.
func GetUserLoader(ctx context.Context) *UserLoader {
    return ctx.Value(userloaderKey).(*UserLoader)
}
