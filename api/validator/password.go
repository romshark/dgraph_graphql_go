package validator

import "github.com/pkg/errors"

// Password implements the Validator interface
func (vld *validator) Password(v string) error {
	if uint(len(v)) < vld.conf.PasswordLenMin {
		return errors.Errorf(
			"password too short (%d / %d)",
			len(v),
			vld.conf.PasswordLenMin,
		)
	}
	if uint(len(v)) > vld.conf.PasswordLenMax {
		return errors.Errorf(
			"password too long (%d / %d)",
			len(v),
			vld.conf.PasswordLenMax,
		)
	}
	return nil
}
