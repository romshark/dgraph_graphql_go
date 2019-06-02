package gqlshield

import (
	"strings"

	uuid "github.com/satori/go.uuid"
)

// ID represents a universally unique identifier
type ID string

// newID creates a new random universally unique identifier
func newID() ID {
	return ID(strings.Replace(uuid.NewV4().String(), "-", "", -1))
}
