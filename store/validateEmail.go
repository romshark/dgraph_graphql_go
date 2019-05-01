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
	if !regexpEmail.MatchString(v) {
		return errors.Errorf("invalid email")
	}
	return nil
}
