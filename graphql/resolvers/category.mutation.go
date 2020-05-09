package resolvers

import (
    "context"

    "github.com/secmohammed/meetups/middlewares"
    "github.com/secmohammed/meetups/models"
    "github.com/secmohammed/meetups/utils/errors"
)

func (c *mutationResolver) CreateInterest(ctx context.Context, categoryID string) (bool, error) {
    currentUser, _ := middlewares.GetCurrentUserFromContext(ctx)

    _, err := c.CategoriesRepo.GetByID(categoryID)
    if err != nil {
        return false, errors.ErrRecordNotFound
    }

    interest := &models.CategoryUser{
        CategoryID: categoryID,
        UserID:     currentUser.ID,
    }
    return c.CategoriesRepo.CreateInterest(interest)
}
func (c *mutationResolver) DeleteInterest(ctx context.Context, categoryID string) (bool, error) {
    currentUser, _ := middlewares.GetCurrentUserFromContext(ctx)

    _, err := c.CategoriesRepo.GetByID(categoryID)
    if err != nil {
        return false, errors.ErrRecordNotFound
    }

    interest := &models.CategoryUser{
        CategoryID: categoryID,
        UserID:     currentUser.ID,
    }
    return c.CategoriesRepo.DeleteInterest(interest)

}

func (c *mutationResolver) CreateCategory(ctx context.Context, input models.CreateCategoryInput) (*models.Category, error) {
    currentUser, _ := middlewares.GetCurrentUserFromContext(ctx)
    if err := input.Validate(); err != nil {
        return nil, err
    }
    category := &models.Category{
        Name:   input.Name,
        UserID: currentUser.ID,
    }
    return c.CategoriesRepo.Create(category)
}
func (c *mutationResolver) UpdateCategory(ctx context.Context, name string, input *models.CreateCategoryInput) (*models.Category, error) {

    category, err := c.CategoriesRepo.GetByName(name)
    if err != nil {
        return nil, errors.ErrRecordNotFound
    }
    if err := input.Validate(); err != nil {
        return nil, err
    }
    category.Name = input.Name
    return c.CategoriesRepo.Update(category)
}
func (c *mutationResolver) DeleteCategory(ctx context.Context, name string) (bool, error) {

    category, err := c.CategoriesRepo.GetByName(name)
    if err != nil {
        return false, errors.ErrRecordNotFound
    }
    return true, c.CategoriesRepo.Delete(category)
}
