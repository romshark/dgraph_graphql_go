package validator

import "github.com/pkg/errors"

// Email implements the Validator interface
func (vld *validator) Email(v string) error {
	if uint(len(v)) > vld.conf.EmailLenMax {
		return errors.Errorf(
			"email address too long (%d / %d)",
			len(v),
			vld.conf.EmailLenMax,
		)
	}
	if !vld.regexpEmail.MatchString(v) {
		return errors.Errorf("invalid email address")
	}
	return nil
}
