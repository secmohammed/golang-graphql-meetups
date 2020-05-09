package middlewares

import (
    "context"
    "net/http"
    "os"
    "strings"

    "github.com/dgrijalva/jwt-go"
    "github.com/dgrijalva/jwt-go/request"
    "github.com/pkg/errors"

    "github.com/secmohammed/meetups/models"
    "github.com/secmohammed/meetups/postgres"
    "github.com/secmohammed/meetups/utils"
)

//CurrentUserKey is used to get the context key of current user.
const CurrentUserKey = utils.ContextKey("currentUser")

//AuthMiddleware is used to authenticate the user via the passed bearer token.
func AuthMiddleware(repo postgres.UsersRepo) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            token, err := parseToken(r)
            if err != nil {
                next.ServeHTTP(w, r)
                return
            }

            claims, ok := token.Claims.(jwt.MapClaims)

            if !ok || !token.Valid {
                next.ServeHTTP(w, r)
                return
            }

            user, err := repo.GetByID(claims["jti"].(string))
            if err != nil {
                next.ServeHTTP(w, r)
                return
            }

            ctx := context.WithValue(r.Context(), CurrentUserKey, user)

            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}

var authHeaderExtractor = &request.PostExtractionFilter{
    Extractor: request.HeaderExtractor{"Authorization"},
    Filter:    stripBearerPrefixFromToken,
}

func stripBearerPrefixFromToken(token string) (string, error) {
    bearer := "BEARER"

    if len(token) > len(bearer) && strings.ToUpper(token[0:len(bearer)]) == bearer {
        return token[len(bearer)+1:], nil
    }

    return token, nil
}

var authExtractor = &request.MultiExtractor{
    authHeaderExtractor,
    request.ArgumentExtractor{"access_token"},
}

func parseToken(r *http.Request) (*jwt.Token, error) {
    jwtToken, err := request.ParseFromRequest(r, authExtractor, func(token *jwt.Token) (interface{}, error) {
        t := []byte(os.Getenv("JWT_SECRET"))
        return t, nil
    })

    return jwtToken, errors.Wrap(err, "parseToken error: ")
}
func HasRole(ctx context.Context, role string, userRepo postgres.UsersRepo) (bool, error) {
    user, err := GetCurrentUserFromContext(ctx)
    if err != nil {
        return false, err
    }
    return userRepo.Can(user.ID, role)

}

//GetCurrentUserFromContext is used to parse the user that we have at our context.
func GetCurrentUserFromContext(ctx context.Context) (*models.User, error) {
    errNoUserInContext := errors.New("no user in context")

    if ctx.Value(CurrentUserKey) == nil {
        return nil, errNoUserInContext
    }

    user, ok := ctx.Value(CurrentUserKey).(*models.User)
    if !ok || user.ID == "" {
        return nil, errNoUserInContext
    }

    return user, nil
}
