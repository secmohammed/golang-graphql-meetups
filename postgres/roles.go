package postgres

import (
    "github.com/go-pg/pg"
    "github.com/secmohammed/meetups/models"
)

// RolesRepo is used to contain the db driver.
type RolesRepo struct {
    DB *pg.DB
}

//Create is used to create a attendee using the passed struct.
func (r *RolesRepo) Create(role *models.Role) (*models.Role, error) {
    _, err := r.DB.Model(role).Returning("*").Insert()
    return role, err
}

// GetBySlug function is used to get the user by its passed slug.
func (r *RolesRepo) GetBySlug(slug string) (*models.Role, error) {
    var role models.Role
    err := r.DB.Model(&role).Where("slug = ?", slug).First()
    if err != nil {
        return nil, err
    }
    return &role, nil
}
