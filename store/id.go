package store

import (
	"errors"
	"regexp"
	"strings"

	uuid "github.com/satori/go.uuid"
)

// ID represents a universally unique identifier
type ID string

// NewID creates a new random universally unique identifier
func NewID() ID {
	return ID(strings.Replace(uuid.NewV4().String(), "-", "", -1))
}

// Verify returns true if the given string represents a valid unique identifier,
// otherwise returns false
func Verify(id string) error {
	res, err := regexp.MatchString("[0-9a-fA-F]{32}", id)
	if err != nil {
		return err
	}
	if !res {
		return errors.New("invalid Identifier value")
	}
	return nil
}

// ImplementsGraphQLType implements the GraphQL scalar type interface
func (id ID) ImplementsGraphQLType(name string) bool {
	return name == "Identifier"
}

// UnmarshalGraphQL implements the GraphQL scalar type interface
func (id *ID) UnmarshalGraphQL(input interface{}) (err error) {
	switch input := input.(type) {
	case string:
		return Verify(input)
	default:
		return errors.New("wrong type, identifier string expected")
	}
}
