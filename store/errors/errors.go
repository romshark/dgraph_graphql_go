package errors

import "fmt"

// Code represents an error code
type Code string

const (
	// ErrInvalidInput is thrown when the API user input is invalid
	ErrInvalidInput Code = "InvalidInput"
)

// Error represents a typed store error
type Error struct {
	Code    string
	Message string
	err     error
}

// filterCode turns unknown error codes to empty strings
func filterCode(code Code) string {
	switch code {
	case ErrInvalidInput:
		return string(code)
	}
	return ""
}

// New creates a new store error
func New(code Code, message string) Error {
	return Error{
		Code:    filterCode(code),
		Message: message,
	}
}

// Newf creates a new store error
func Newf(code Code, format string, v ...interface{}) Error {
	return Error{
		Code:    filterCode(code),
		Message: fmt.Sprintf(format, v...),
	}
}

// Wrap wraps an error into a new store error
func Wrap(code Code, err error) Error {
	return Error{
		Code:    filterCode(code),
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
