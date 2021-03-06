package errors

import "errors"

var (
    //ErrBadCredentials is used to clarify that user has given invaild credentials.
    ErrBadCredentials = errors.New("Invalid credentials")
    //ErrAuthenticated is used to clarify that user is already authenticated.
    ErrAuthenticated = errors.New("you are already authenticated")
    // ErrUnauthenticated is used to indicate that user is unauthenticated.
    ErrUnauthenticated = errors.New("Unauthorized Attempt")
    //ErrRecordNotFound is used to indicate that there is no record at database.
    ErrRecordNotFound = errors.New("record doesn't exist")
    //ErrInternalError is used to indicate that there is something went wrong with the server.
    ErrInternalError = errors.New("Whoops, something went wrong")
    //ErrInsufficientPermissions is used to indicate that user doesn't have a certain permission.
    ErrInsufficientPermissions = errors.New("You don't have this permission to proceed")
    //ErrEmailIsntUnique is used to indicate that email isn't unique at database.
    ErrEmailIsntUnique = errors.New("email is already used")
    //ErrUsernameIsntUnique  is used to indicate that email isn't unique at database.
    ErrUsernameIsntUnique = errors.New("username is already used")
    //ErrCouldntGenerateJWTToken is used to indicate that here was an error while generating the jwt token.
    ErrCouldntGenerateJWTToken = errors.New("failed generating token")
    //ErrCouldntAssignMemberToGroup is used to show an error message in case of unauthorized attempt while adding user to group.
    ErrCouldntAssignMemberToGroup = errors.New("couldn't assign user to group, due you aren't allowed to add other members.")
)
