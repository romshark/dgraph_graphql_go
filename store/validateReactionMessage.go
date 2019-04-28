package store

import "github.com/pkg/errors"

// ValidateReactionMessage returns an error if invalid, otherwise returns nil
func ValidateReactionMessage(v string) error {
	if len(v) < 1 {
		return errors.Errorf("Reaction.message too short (min: 1)")
	}
	if len(v) > 256 {
		return errors.Errorf("Reaction.message too long (%d / 256)", len(v))
	}
	return nil
}
