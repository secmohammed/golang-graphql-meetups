package postgres

import (
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
    query := c.DB.Model(&categories).Order("id")
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
    err := c.DB.Model(&category).Where("name = ?", name).First()
    if err != nil {
        return nil, err
    }
    return &category, nil
}

// GetByID is used to fetch category by id.
func (c *CategoriesRepo) GetByID(id string) (*models.Category, error) {
    category := models.Category{}
    err := c.DB.Model(&category).Where("id = ?", id).Select()
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

//CreateInterest is used to create an interest using the passed struct.
func (c *CategoriesRepo) CreateInterest(interest *models.CategoryUser) (bool, error) {
    _, err := c.DB.Model(interest).Returning("*").Insert()
    if err != nil {
        return false, err
    }
    return true, err
}

//DeleteInterest is used to remove an interest of a user.
func (c *CategoriesRepo) DeleteInterest(interest *models.CategoryUser) (bool, error) {
    _, err := c.DB.Model(interest).Where("user_id = ?", interest.UserID).Where("category_id = ?", interest.CategoryID).Delete()
    if err != nil {
        return false, err
    }
    return true, err
}
