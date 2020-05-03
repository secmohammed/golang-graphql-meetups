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
    UserByID            *UserLoader
    MeetupsByCategory   *MeetupsLoader
    MeetupsByGroup      *MeetupsLoader
    CategoriesByMeetup  *CategoriesLoader
    CategoriesByGroup   *CategoriesLoader
    MembersByGroup      *MembersLoader
    AttendeesByMeetup   *AttendeesLoader
    NotificationsByUser *NotificationsLoader
    MeetupByID          *MeetupLoader
    InterestsByUser     *CategoriesLoader
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
        loaders.MeetupsByGroup = &MeetupsLoader{
            wait:     wait,
            maxBatch: 100,
            fetch: func(ids []string) ([][]*models.Meetup, []error) {
                var groups []*models.Group
                err := db.Model(&groups).Where("id in (?)", pg.In(ids)).OrderExpr("id DESC").Relation("Meetups").Select()
                if err != nil {
                    return nil, []error{err}
                }
                c := make(map[string]*models.Group, len(ids))
                for _, group := range groups {
                    c[group.ID] = group
                }
                meetups := make([][]*models.Meetup, len(ids))
                for i, id := range ids {
                    meetups[i] = c[id].Meetups
                }
                return meetups, nil
            },
        }
        loaders.MembersByGroup = &MembersLoader{
            wait:     wait,
            maxBatch: 100,
            fetch: func(ids []string) ([][]*models.UserGroup, []error) {
                var groups []*models.Group
                err := db.Model(&groups).Where("\"group\".\"id\" in (?)", pg.In(ids)).OrderExpr("id DESC").Relation("Members").Select()
                if err != nil {
                    return nil, []error{err}
                }
                c := make(map[string]*models.Group, len(ids))
                for _, group := range groups {
                    c[group.ID] = group
                }
                members := make([][]*models.UserGroup, len(ids))
                for i, id := range ids { // group
                    for _, member := range c[id].Members {
                        members[i] = append(members[i], &models.UserGroup{User: member, Type: member.Type})
                    }
                }
                return members, nil
            },
        }
        loaders.NotificationsByUser = &NotificationsLoader{
            wait:     wait,
            maxBatch: 100,
            fetch: func(ids []string) ([][]*models.Notification, []error) {
                var users []*models.User
                err := db.Model(&users).Where("id in (?)", pg.In(ids)).OrderExpr("id DESC").Relation("Notifications").Select()
                if err != nil {
                    return nil, []error{err}
                }
                m := make(map[string]*models.User, len(ids))
                for _, user := range users {
                    m[user.ID] = user
                }
                notifications := make([][]*models.Notification, len(ids))
                for i, id := range ids {
                    notifications[i] = m[id].Notifications
                }
                return notifications, nil
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
        loaders.CategoriesByGroup = &CategoriesLoader{
            wait:     wait,
            maxBatch: 100,
            fetch: func(ids []string) ([][]*models.Category, []error) {
                var groups []*models.Group
                err := db.Model(&groups).Where("id in (?)", pg.In(ids)).OrderExpr("id DESC").Relation("Categories").Select()
                if err != nil {
                    return nil, []error{err}
                }
                m := make(map[string]*models.Group, len(ids))
                for _, group := range groups {
                    m[group.ID] = group
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
