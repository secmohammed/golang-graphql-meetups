package models

import (
    "os"
    "time"

    "github.com/99designs/gqlgen/graphql"
    "github.com/dgrijalva/jwt-go"
    "github.com/go-playground/validator"

    "golang.org/x/crypto/bcrypt"
)

//User model attributes.
type User struct {
    ID         string `json:"id"`
    Username   string `json:"username"`
    Email      string `json:"email"`
    Password   string `json:"password"`
    FirstName  string `json:"first_name"`
    LastName   string `json:"last_name"`
    Attendees  []*Attendee
    Categories []*Category `pg:"many2many:category_user"`
    Groups     []*Group    `pg:"many2many:group_user"`
    Avatar     string      `json:"avatar"`
    Type       string      `sql:"-" json:"-"`

    CreatedAt time.Time  `json:"created_at"`
    UpdatedAt time.Time  `json:"updated_at"`
    DeletedAt *time.Time `json:"-" pg:",soft_delete"`
}

// CategoryUser struct type
type CategoryUser struct {
    tableName struct{} `sql:"category_user"`

    CategoryID string `json:"category_id"`
    UserID     string `json:"user_id"`
}

//RegisterInput is used to validate the user against passed inputs while registration.
type RegisterInput struct {
    Username             string          `json:"username" validate:"required,min=3,max=32"`
    Email                string          `json:"email" validate:"required,email"`
    Password             string          `json:"password" validate:"required,min=8,max=32,eqfield=PasswordConfirmation"`
    PasswordConfirmation string          `json:"password_confirmation" validate:"required"`
    FirstName            string          `json:"first_name" validate:"required,min=3,max=32"`
    LastName             string          `json:"last_name" validate:"required,min=3,max=32"`
    Avatar               *graphql.Upload `json:"avatar"`
}

//LoginInput is used to validate the user against passed inputs while logging in.
type LoginInput struct {
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=8,max=32"`
}

//Validate is used to validate teh passed values against the login struct.
func (u *LoginInput) Validate() error {
    validate := validator.New()
    return validate.Struct(u)

}

//Validate is used to validate the passed values against the struct validation props.
func (u *RegisterInput) Validate() error {
    validate := validator.New()
    // Custom rule for validating against column.

    // validate.RegisterValidation("unique", func(fl validator.FieldLevel) bool {
    //     param := strings.Split(fl.Param(), `:`)
    //     paramField := param[0]
    //     paramValue := param[1]
    //     user, err := postgres.GetByField(paramField, paramValue)
    //     if err != nil || user != nil {
    //         return false
    //     }
    //     return true
    // })
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

// ComparePassword is used to compare the plaintext password that user passed against its hashed password which stored in database.
func (u *User) ComparePassword(password string) error {
    bytePassword := []byte(password)
    byteHashedPassword := []byte(u.Password)
    return bcrypt.CompareHashAndPassword(byteHashedPassword, bytePassword)
}
