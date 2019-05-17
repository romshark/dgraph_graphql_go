package validator

import "github.com/pkg/errors"

// Password implements the Validator interface
func (vld *validator) Password(v string) error {
	if uint(len(v)) < vld.opts.PasswordLenMin {
		return errors.Errorf(
			"password too short (%d / %d)",
			len(v),
			vld.opts.PasswordLenMin,
		)
	}
	if uint(len(v)) > vld.opts.PasswordLenMax {
		return errors.Errorf(
			"password too long (%d / %d)",
			len(v),
			vld.opts.PasswordLenMax,
		)
	}
	return nil
}
