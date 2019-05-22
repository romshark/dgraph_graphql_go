package errors

import "fmt"

// Code represents an error code
type Code string

const (
	// ErrUnauthorized is thrown when the API user is unauthorized to perform
	// the request
	ErrUnauthorized Code = "Unauthorized"

	// ErrInvalidInput is thrown when the API user input is invalid
	ErrInvalidInput Code = "InvalidInput"

	// ErrWrongCreds is thrown when the API user provides wrong authentication
	// credentials
	ErrWrongCreds Code = "WrongCreds"
)

// Error represents a typed store error
type Error struct {
	Code    string
	Message string
}

// FilterCode turns unknown error codes to empty strings
func FilterCode(code Code) string {
	switch code {
	case ErrUnauthorized:
		return string(code)
	case ErrInvalidInput:
		return string(code)
	case ErrWrongCreds:
		return string(code)
	}
	return ""
}

// New creates a new store error
func New(code Code, message string) Error {
	return Error{
		Code:    FilterCode(code),
		Message: message,
	}
}

// Newf creates a new store error
func Newf(code Code, format string, v ...interface{}) Error {
	return Error{
		Code:    FilterCode(code),
		Message: fmt.Sprintf(format, v...),
	}
}

// Wrap wraps an error into a new store error
func Wrap(code Code, err error) Error {
	return Error{
		Code:    FilterCode(code),
		Message: err.Error(),
	}
}

// Error implements the error interface
func (err Error) Error() string {
	return err.Message
}

// ErrorCode returns the code if the given error is a store error and has a code
// assigned, otherwise returns an empty string
func ErrorCode(err error) string {
	r, ok := err.(Error)
	if !ok {
		return ""
	}
	return r.Code
}
