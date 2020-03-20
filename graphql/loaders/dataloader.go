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

//Loaders struct.
type Loaders struct {
    UserByID *UserLoader
}

//DataloaderMiddleware is used to load the data related to users.
func DataloaderMiddleware(db *pg.DB, next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        loaders := Loaders{}
        wait := 1 * time.Millisecond
        loaders.UserByID = &UserLoader{
            maxBatch: 100,
            wait:     wait,
            fetch: func(ids []string) ([]*models.User, []error) {
                var users []*models.User
                err := db.Model(&users).Where("id in (?)", pg.In(ids)).OrderExpr("id DESC").Select()
                if err != nil {
                    return nil, []error{err}
                }
                u := make(map[string]*models.User, len(ids))
                for _, user := range users {
                    u[user.ID] = user
                }
                result := make([]*models.User, len(ids))
                for i, id := range ids {
                    result[i] = u[id]
                }
                return result, nil
            },
        }

        ctx := context.WithValue(r.Context(), userloaderKey, loaders)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

//GetLoaders is  fetch the  loaders.
func GetLoaders(ctx context.Context) Loaders {
    return ctx.Value(userloaderKey).(Loaders)
}
