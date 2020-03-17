package resolvers

import (
    "context"
    "errors"
    "log"

    "github.com/secmohammed/meetups/models"
)

var (
    //ErrBadCredentials is used to clarify that user has given invaild credentials.
    ErrBadCredentials = errors.New("Invalid credentials")
)

func (m *mutationResolver) Login(ctx context.Context, input *models.LoginInput) (*models.Auth, error) {
    if err := input.Validate(); err != nil {
        return nil, err
    }

    user, err := m.UsersRepo.GetByField("email", input.Email)
    if err != nil {
        return nil, ErrBadCredentials
    }
    err = user.ComparePassword(input.Password)
    if err != nil {
        return nil, ErrBadCredentials
    }
    token, err := user.GenerateToken()
    if err != nil {
        return nil, errors.New("something went wrong")
    }
    return &models.Auth{
        AuthToken: token,
        User:      user,
    }, nil

}

// Register is used to create a user by the passed attribtues.
func (m *mutationResolver) Register(ctx context.Context, input *models.RegisterInput) (*models.Auth, error) {

    if err := input.Validate(); err != nil {
        return nil, err
    }
    _, err := m.UsersRepo.GetByField("email", input.Email)
    if err == nil {
        return nil, errors.New("email already in used")
    }

    _, err = m.UsersRepo.GetByField("username", input.Username)
    if err == nil {
        return nil, errors.New("username already in used")
    }
    user := &models.User{
        Username:  input.Username,
        Email:     input.Email,
        Password:  input.Password,
        FirstName: input.FirstName,
        LastName:  input.LastName,
    }
    err = user.HashPassword(input.Password)
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
