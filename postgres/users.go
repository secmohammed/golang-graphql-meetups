package postgres

import (
    "github.com/go-pg/pg"
    "github.com/secmohammed/meetups/models"
)

// UsersRepo is used to contain the db driver.
type UsersRepo struct {
    DB *pg.DB
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
