package gqlshield

import "fmt"

// ErrorCode represents an error code
type ErrorCode string

const (
	// ErrUnauthorized is returned when Check denies access
	ErrUnauthorized ErrorCode = "Unauthorized"

	// ErrWrongInput is returned when Check fails due to a client error
	ErrWrongInput ErrorCode = "WrongInput"
)

// Error represents a typed GraphQL shield error
type Error struct {
	Code    ErrorCode
	Message string
}

func (err Error) Error() string {
	return fmt.Sprintf("%s (%s)", err.Message, err.Code)
}

// ErrCode returns the code of the GraphQLShield error (if it is one)
func ErrCode(err error) ErrorCode {
	if typedErr, isError := err.(Error); isError {
		return typedErr.Code
	}
	return ""
}
