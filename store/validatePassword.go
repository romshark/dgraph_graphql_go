package store

import (
	"github.com/pkg/errors"
)

// ValidatePassword returns an error if invalid, otherwise returns nil
func ValidatePassword(v string) error {
	if len(v) < 6 {
		return errors.Errorf("password too short (%d / 6)", len(v))
	}
	if len(v) > 256 {
		return errors.Errorf("password too long (%d / 256)", len(v))
	}
	return nil
}
