package resolvers

import (
    "context"
    "errors"

    "github.com/secmohammed/meetups/middlewares"
    "github.com/secmohammed/meetups/models"
)

func (c *mutationResolver) CreateCategory(ctx context.Context, input models.CreateCategoryInput) (*models.Category, error) {
    currentUser, err := middlewares.GetCurrentUserFromContext(ctx)
    if err != nil {
        return nil, ErrUnauthenticated
    }
    if err := input.Validate(); err != nil {
        return nil, err
    }
    // TODO: check if user is admin, or has the ability to create a category.
    category := &models.Category{
        Name:   input.Name,
        UserID: currentUser.ID,
    }
    return c.CategoriesRepo.Create(category)
}
func (c *mutationResolver) UpdateCategory(ctx context.Context, name string, input *models.CreateCategoryInput) (*models.Category, error) {
    currentUser, err := middlewares.GetCurrentUserFromContext(ctx)
    if err != nil {
        return nil, ErrUnauthenticated
    }

    category, err := c.CategoriesRepo.GetByName(name)
    if err != nil {
        return nil, errors.New("Couldn't find this category to update")
    }
    if category.UserID != currentUser.ID {
        return nil, errors.New("Unauthorized attempt")
    }
    if err := input.Validate(); err != nil {
        return nil, err
    }
    category.Name = input.Name
    return c.CategoriesRepo.Update(category)
}
func (c *mutationResolver) DeleteCategory(ctx context.Context, name string) (bool, error) {
    currentUser, err := middlewares.GetCurrentUserFromContext(ctx)
    if err != nil {
        return false, ErrUnauthenticated
    }

    category, err := c.CategoriesRepo.GetByName(name)
    if err != nil {
        return false, errors.New("Couldn't find this category to update")
    }
    if category.UserID != currentUser.ID {
        return false, errors.New("Unauthorized attempt")
    }

    return true, c.CategoriesRepo.Delete(category)
}
