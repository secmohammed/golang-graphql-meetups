package postgres

import (
    "fmt"

    "github.com/go-pg/pg"
    "github.com/secmohammed/meetups/models"
)

// UsersRepo is used to contain the db driver.
type UsersRepo struct {
    DB *pg.DB
}

//GetByField is used to retrieve a user model by a specific criteria.
func (u *UsersRepo) GetByField(field, value string) (*models.User, error) {
    var user models.User
    err := u.DB.Model(&user).Where(fmt.Sprintf("%v = ?", field), value).First()
    return &user, err
}

// GetByID function is used to get the user by its passed id.
func (u *UsersRepo) GetByID(id string) (*models.User, error) {
    var user models.User
    err := u.DB.Model(&user).Where("id = ?", id).First()
    if err != nil {
        return nil, err
    }
    return &user, nil
}

// CreateUser is used to create user for the database.
func (u *UsersRepo) CreateUser(tx *pg.Tx, user *models.User) (*models.User, error) {
    _, err := tx.Model(user).Returning("*").Insert()
    return user, err
}
