package postgres

import (
    "fmt"

    "github.com/go-pg/pg"
    "github.com/secmohammed/meetups/models"
    "github.com/secmohammed/meetups/utils/errors"
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
func (u *UsersRepo) Can(id, role string) (bool, error) {
    // fetch the user by id
    // load the roles for this user
    var user models.User
    err := u.DB.Model(&user).Where("id = ?", id).Relation("Roles").First()
    if err != nil {
        return false, errors.ErrRecordNotFound
    }
    // merge his permissions with the roles permissions
    var permissions map[string]bool
    for _, role := range user.Roles {
        for key, value := range role.Permissions {
            permissions[key] = value
        }
    }
    for key, value := range user.Permissions {
        permissions[key] = value
    }
    // check if the map contains the passed role
    if permissions[role] {
        return true, nil
    }
    return false, errors.ErrInternalError
}
func (u *UsersRepo) Cannot(id, role string) (bool, error) {
    ok, err := u.Can(id, role)
    if err != nil {
        return false, err
    }
    return !ok, nil
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
func (u *UsersRepo) AttachUserToRole(user *models.User, role *models.Role) (bool, error) {
    userRole := &models.RoleUser{
        UserID: user.ID,
        RoleID: role.ID,
    }
    _, err := u.DB.Model(&userRole).Insert()
    if err != nil {
        return false, err
    }
    return true, nil

}

// CreateUser is used to create user for the database.
func (u *UsersRepo) CreateUser(tx *pg.Tx, user *models.User) (*models.User, error) {
    _, err := tx.Model(user).Returning("*").Insert()
    return user, err
}
