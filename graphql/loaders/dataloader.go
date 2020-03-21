package loaders

import (
    "context"
    "net/http"
    "time"

    "github.com/go-pg/pg"
    "github.com/go-pg/pg/orm"
    "github.com/secmohammed/meetups/models"
    "github.com/secmohammed/meetups/utils"
)

const userloaderKey = utils.ContextKey("userloader")

//Loaders struct.
type Loaders struct {
    UserByID           *UserLoader
    MeetupsByCategory  *MeetupsLoader
    CategoriesByMeetup *CategoriesLoader
    AttendeesByMeetup  *AttendeesLoader
    MeetupByID         *MeetupLoader
    InterestsByUser    *CategoriesLoader
    // commentsByMeetup    *ItemSliceLoader
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
        loaders.MeetupByID = &MeetupLoader{
            maxBatch: 100,
            wait:     wait,
            fetch: func(ids []string) ([]*models.Meetup, []error) {
                var meetups []*models.Meetup
                err := db.Model(&meetups).Where("id in (?)", pg.In(ids)).OrderExpr("id DESC").Select()
                if err != nil {
                    return nil, []error{err}
                }
                m := make(map[string]*models.Meetup, len(ids))
                for _, meetup := range meetups {
                    m[meetup.ID] = meetup
                }
                result := make([]*models.Meetup, len(ids))
                for i, id := range ids {
                    result[i] = m[id]
                }
                return result, nil
            },
        }

        loaders.MeetupsByCategory = &MeetupsLoader{
            wait:     wait,
            maxBatch: 100,
            fetch: func(ids []string) ([][]*models.Meetup, []error) {
                var categories []*models.Category
                err := db.Model(&categories).Where("id in (?)", pg.In(ids)).OrderExpr("id DESC").Relation("Meetups").Select()
                if err != nil {
                    return nil, []error{err}
                }
                c := make(map[string]*models.Category, len(ids))
                for _, category := range categories {
                    c[category.ID] = category
                }
                meetups := make([][]*models.Meetup, len(ids))
                for i, id := range ids {
                    meetups[i] = c[id].Meetups
                }
                return meetups, nil
            },
        }
        loaders.CategoriesByMeetup = &CategoriesLoader{
            wait:     wait,
            maxBatch: 100,
            fetch: func(ids []string) ([][]*models.Category, []error) {
                var meetups []*models.Meetup
                err := db.Model(&meetups).Where("id in (?)", pg.In(ids)).OrderExpr("id DESC").Relation("Categories").Select()
                if err != nil {
                    return nil, []error{err}
                }
                m := make(map[string]*models.Meetup, len(ids))
                for _, meetup := range meetups {
                    m[meetup.ID] = meetup
                }
                categories := make([][]*models.Category, len(ids))
                for i, id := range ids {
                    categories[i] = m[id].Categories
                }
                return categories, nil
            },
        }
        loaders.InterestsByUser = &CategoriesLoader{
            wait:     wait,
            maxBatch: 100,
            fetch: func(ids []string) ([][]*models.Category, []error) {
                var categories []*models.Category
                err := db.Model(&categories).Relation("Users", func(q *orm.Query) (*orm.Query, error) {
                    return q.Where("category_user.user_id in (?)", pg.In(ids)), nil
                }).Select()

                if err != nil {
                    return nil, []error{err}
                }
                categoriesCollection := make([][]*models.Category, len(ids))

                for _, category := range categories {
                    for _, user := range category.Users {
                        for i, id := range ids {
                            if id == user.ID {
                                categoriesCollection[i] = append(categoriesCollection[i], category) 
                            }
                        }
                    }
                }
                return categoriesCollection, nil
            },
        }
        loaders.AttendeesByMeetup = &AttendeesLoader{
            wait:     wait,
            maxBatch: 100,
            fetch: func(ids []string) ([][]*models.Attendee, []error) {
                var meetups []*models.Meetup
                err := db.Model(&meetups).Where("id in (?)", pg.In(ids)).OrderExpr("id DESC").Relation("Attendees").Select()
                if err != nil {
                    return nil, []error{err}
                }
                m := make(map[string]*models.Meetup, len(ids))
                for _, meetup := range meetups {
                    m[meetup.ID] = meetup
                }
                attendees := make([][]*models.Attendee, len(ids))
                for i, id := range ids {
                    attendees[i] = m[id].Attendees
                }
                return attendees, nil
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
