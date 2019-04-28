package store

import "github.com/pkg/errors"

// ValidateUserDisplayName returns an error if invalid, otherwise returns nil
func ValidateUserDisplayName(v string) error {
	if len(v) < 2 {
		return errors.Errorf("User.displayName too short (min: 2)")
	}
	if len(v) > 64 {
		return errors.Errorf("User.displayName too long (%d / 64)", len(v))
	}
	return nil
}
