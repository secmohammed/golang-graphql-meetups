package postgres

import (
    "github.com/go-pg/pg"
    "github.com/secmohammed/meetups/models"
)

// GroupsRepo is used to contain the db driver.
type GroupsRepo struct {
    DB *pg.DB
}

//Create is used to create a group using the passed struct.
func (g *GroupsRepo) Create(group *models.Group) (*models.Group, error) {
    _, err := g.DB.Model(group).Returning("*").Insert()
    return group, err
}

// GetByID is used to fetch meetup by id.
func (g *GroupsRepo) GetByID(id string) (*models.Group, error) {
    group := models.Group{}
    err := g.DB.Model(&group).Where("id = ?", id).First()
    if err != nil {
        return nil, err
    }
    return &group, nil
}

// GetGroups is used to get meetups from database.
func (g *GroupsRepo) GetGroups() ([]*models.Group, error) {
    var groups []*models.Group
    err := g.DB.Model(&groups).Order("id").Select()
    if err != nil {
        return nil, err
    }
    return groups, nil
}

// Update is used to update the passed group by id.
func (g *GroupsRepo) Update(group *models.Group) (*models.Group, error) {
    _, err := g.DB.Model(group).Where("id = ?", group.ID).Update()
    return group, err
}

// Delete is used to delete group by its id.
func (g *GroupsRepo) Delete(group *models.Group) error {
    _, err := g.DB.Model(group).Where("id = ?", group.ID).Delete()
    return err
}
