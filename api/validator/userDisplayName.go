package validator

import "github.com/pkg/errors"

// UserDisplayName implements the Validator interface
func (vld *validator) UserDisplayName(v string) error {
	if uint(len(v)) < vld.opts.UserDisplayNameLenMin {
		return errors.Errorf(
			"User.displayName too short (min: %d)",
			vld.opts.UserDisplayNameLenMin,
		)
	}
	if uint(len(v)) > vld.opts.UserDisplayNameLenMax {
		return errors.Errorf(
			"User.displayName too long (%d / %d)",
			len(v),
			vld.opts.UserDisplayNameLenMax,
		)
	}
	return nil
}
