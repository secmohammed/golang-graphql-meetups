package models

import (
    "os"
    "time"

    "github.com/dgrijalva/jwt-go"

    "github.com/go-playground/validator"
    "golang.org/x/crypto/bcrypt"
)

//User model attributes.
type User struct {
    ID        string     `json:"id"`
    Username  string     `json:"username"`
    Email     string     `json:"email"`
    Password  string     `json:"password"`
    FirstName string     `json:"first_name"`
    LastName  string     `json:"last_name"`
    CreatedAt time.Time  `json:"created_at"`
    UpdatedAt time.Time  `json:"updated_at"`
    DeletedAt *time.Time `json:"-" pg:",soft_delete"`
}

//Validate is used to validate the passed values against the struct validation props.
func (u *User) Validate() error {
    validate := validator.New()
    return validate.Struct(u)
}

//HashPassword is used to hash the passed password.
func (u *User) HashPassword(password string) error {
    bytePassword := []byte(password)
    passwordHash, err := bcrypt.GenerateFromPassword(bytePassword, bcrypt.DefaultCost)
    if err != nil {
        return err
    }
    u.Password = string(passwordHash)

    return nil
}

//GenerateToken is used to generate the bearer token authorization.
func (u *User) GenerateToken() (*AuthToken, error) {
    expiredAt := time.Now().Add(time.Hour * 24 * 7) // a week

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
        ExpiresAt: expiredAt.Unix(),
        Id:        u.ID,
        IssuedAt:  time.Now().Unix(),
        Issuer:    "meetmeup",
    })

    accessToken, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
    if err != nil {
        return nil, err
    }
    return &AuthToken{
        AccessToken: accessToken,
        ExpiredAt:   expiredAt,
    }, nil
}
