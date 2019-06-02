package gqlshield

import (
	"regexp"
	"strings"

	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

// ID represents a universally unique identifier
type ID string

// newID creates a new random universally unique identifier
func newID() ID {
	return ID(strings.Replace(uuid.NewV4().String(), "-", "", -1))
}

// Validate returns an error if the ID has an invalid value
func (id ID) Validate() error {
	res, err := regexp.MatchString("[0-9a-fA-F]{32}", string(id))
	if err != nil {
		return err
	}
	if !res {
		return errors.New("invalid ID value")
	}
	return nil
}
