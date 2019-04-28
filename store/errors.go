package store

type TypedError struct {
	Code    string
	Message string
}

func ParseError(code string) (TypedError, bool) {
	switch code {
	}
	return TypedError{}, false
}
