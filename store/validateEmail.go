package store

import (
	"regexp"

	"github.com/pkg/errors"
)

var regexpEmail *regexp.Regexp

func init() {
	var err error
	regexpEmail, err = regexp.Compile("^.+@.+\\..+$")
	if err != nil {
		panic(errors.Wrap(err, "compile regexpEmail"))
	}
}

// ValidateEmail returns an error if invalid, otherwise returns nil
func ValidateEmail(v string) error {
	if len(v) > 96 {
		return errors.Errorf("email address too long (%d)", len(v))
	}
	if !regexpEmail.MatchString(v) {
		return errors.Errorf("invalid email address")
	}
	return nil
}
