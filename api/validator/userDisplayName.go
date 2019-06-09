package validator

import "github.com/pkg/errors"

// UserDisplayName implements the Validator interface
func (vld *validator) UserDisplayName(v string) error {
	if uint(len(v)) < vld.conf.UserDisplayNameLenMin {
		return errors.Errorf(
			"User.displayName too short (min: %d)",
			vld.conf.UserDisplayNameLenMin,
		)
	}
	if uint(len(v)) > vld.conf.UserDisplayNameLenMax {
		return errors.Errorf(
			"User.displayName too long (%d / %d)",
			len(v),
			vld.conf.UserDisplayNameLenMax,
		)
	}
	return nil
}
