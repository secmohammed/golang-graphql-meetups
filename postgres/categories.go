package postgres

import (
    "fmt"

    "github.com/go-pg/pg"
    "github.com/secmohammed/meetups/models"
)

// CategoriesRepo is used to contain the db driver.
type CategoriesRepo struct {
    DB *pg.DB
}

// GetCategories is used to get categories from database.
func (c *CategoriesRepo) GetCategories(limit, offset *int) ([]*models.Category, error) {
    var categories []*models.Category
    query := c.DB.Model(&categories).Relation("Meetups").Order("id")
    if limit != nil {
        query.Limit(*limit)
    }
    if offset != nil {
        query.Offset(*offset)
    }
    err := query.Select()
    if err != nil {
        return nil, err
    }
    return categories, nil
}

//Create is used to create a comment using the passed struct.
func (c *CategoriesRepo) Create(category *models.Category) (*models.Category, error) {
    _, err := c.DB.Model(category).Returning("*").Insert()
    return category, err
}

// Update is used to update the passed meetup by id.
func (c *CategoriesRepo) Update(category *models.Category) (*models.Category, error) {
    _, err := c.DB.Model(category).Where("id = ?", category.ID).Update()
    return category, err
}

// GetByName is used to fetch meetup by name.
func (c *CategoriesRepo) GetByName(name string) (*models.Category, error) {
    category := models.Category{}
    err := c.DB.Model(&category).Where("name = ?", name).Relation("Meetups").First()
    if err != nil {
        return nil, err
    }
    return &category, nil
}

// Delete is used to delete meetup by its id.
func (c *CategoriesRepo) Delete(category *models.Category) error {
    _, err := c.DB.Model(category).Where("id = ?", category.ID).Delete()
    return err
}

//GetMeetupsForCategory is used to get the meetups for th associated category.
func (c *CategoriesRepo) GetMeetupsForCategory(id string) ([]*models.Meetup, error) {
    category := new(models.Category)
    err := c.DB.Model(&category).Relation("Meetups").Where("id = ?", id).First()
    meetups := make([]*models.Meetup, len(category.Meetups))
    for _, meetup := range category.Meetups {
        fmt.Println(meetup)
        meetups = append(meetups, meetup)
    }
    return meetups, err
}

// GetCategoriesForUser is used to get comments for the passed user by its id.
func (c *CategoriesRepo) GetCategoriesForUser(user *models.User) ([]*models.Category, error) {
    var categories []*models.Category
    err := c.DB.Model(&categories).Where("user_id = ? ", user.ID).Order("id").Select()
    return categories, err
}
