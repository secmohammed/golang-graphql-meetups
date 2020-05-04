package resolvers

import (
    "context"
    coreErrors "errors"
    "log"
    "os"

    "github.com/secmohammed/meetups/models"
    "github.com/secmohammed/meetups/utils"
    "github.com/secmohammed/meetups/utils/errors"
)

func (m *mutationResolver) Login(ctx context.Context, input *models.LoginInput) (*models.Auth, error) {
    if err := input.Validate(); err != nil {
        return nil, err
    }

    user, err := m.UsersRepo.GetByField("email", input.Email)
    if err != nil {
        log.Println(err)
        return nil, errors.ErrBadCredentials
    }
    err = user.ComparePassword(input.Password)
    if err != nil {
        return nil, errors.ErrBadCredentials
    }
    token, err := user.GenerateToken()
    if err != nil {
        return nil, errors.ErrInternalError
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
        return nil, errors.ErrEmailIsntUnique
    }

    _, err = m.UsersRepo.GetByField("username", input.Username)
    if err == nil {
        return nil, errors.ErrUsernameIsntUnique
    }
    user := &models.User{
        Username:  input.Username,
        Email:     input.Email,
        Password:  input.Password,
        FirstName: input.FirstName,
        LastName:  input.LastName,
    }
    if input.Avatar != nil {
        info, status := utils.UploadFile(input.Avatar, os.Getenv("AVATAR_STORAGE_PATH"))
        if !status {
            return nil, coreErrors.New(info)
        }
        if status {
            user.Avatar = info
        }

    }

    err = user.HashPassword(input.Password)
    if err != nil {
        return nil, errors.ErrInternalError
    }
    transaction, err := m.UsersRepo.DB.Begin()
    if err != nil {
        log.Printf("error creating a transaction: %v", err)
        return nil, errors.ErrInternalError
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
        return nil, errors.ErrCouldntGenerateJWTToken
    }
    return &models.Auth{
        AuthToken: token,
        User:      user,
    }, nil
}
