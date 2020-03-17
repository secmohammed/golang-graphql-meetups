package resolvers

import (
    "context"
    "errors"
    "log"

    "github.com/secmohammed/meetups/models"
)

// Register is used to create a user by the passed attribtues.
func (m *mutationResolver) Register(ctx context.Context, input *models.RegisterInput) (*models.Auth, error) {
    user := &models.User{
        Username:  input.Username,
        Email:     input.Email,
        Password:  input.Password,
        FirstName: input.FirstName,
        LastName:  input.LastName,
    }
    if err := user.Validate(); err != nil {
        return nil, err
    }
    err := user.HashPassword(input.Password)
    if err != nil {
        return nil, errors.New("something went wrong")
    }
    transaction, err := m.UsersRepo.DB.Begin()
    if err != nil {
        log.Printf("error creating a transaction: %v", err)
        return nil, errors.New("something went wrong")
    }
    defer transaction.Rollback()
    if _, err := m.UsersRepo.CreateUser(transaction, user); err != nil {
        log.Printf("error creating a user: %v", err)

        return nil, err
    }
    if err := transaction.Commit(); err != nil {
        log.Printf("error while Commiting: %v", err)
        return nil, err

    }
    token, err := user.GenerateToken()
    if err != nil {
        log.Printf("error while Commiting: %v", err)
        return nil, errors.New("failed generating token")
    }
    return &models.Auth{
        AuthToken: token,
        User:      user,
    }, nil
}
