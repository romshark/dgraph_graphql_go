package validator

import (
	"fmt"
	"regexp"

	"github.com/pkg/errors"
)

// Validator represents a validator interface
type Validator interface {
	// Email returns an error if the given email address is invalid,
	// otherwise returns nil
	Email(v string) error

	// Password returns an error if the given password is invalid,
	// otherwise returns nil
	Password(v string) error

	// PostContents returns an error if the given post contents are invalid,
	// otherwise returns nil
	PostContents(v string) error

	// PostTitle returns an error if the given post title is invalid,
	// otherwise returns nil
	PostTitle(v string) error

	// ReactionMessage returns an error if the given reaction message is
	// invalid, otherwise returns nil
	ReactionMessage(v string) error

	// UserDisplayName returns an error if the given user display name is
	// invalid, otherwise returns nil
	UserDisplayName(v string) error
}

// Config represents the validator configuration
type Config struct {
	PasswordLenMin        uint
	PasswordLenMax        uint
	EmailLenMax           uint
	PostContentsLenMin    uint
	PostContentsLenMax    uint
	PostTitleLenMin       uint
	PostTitleLenMax       uint
	ReactionMessageLenMin uint
	ReactionMessageLenMax uint
	UserDisplayNameLenMin uint
	UserDisplayNameLenMax uint
}

type validator struct {
	conf        Config
	regexpEmail *regexp.Regexp
}

// NewValidator creates a new validator instance
func NewValidator(productionModeEnabled bool, conf Config) (Validator, error) {
	regexpEmail, err := regexp.Compile(`^.+@.+\..+$`)
	if err != nil {
		panic(errors.Wrap(err, "compile regexpEmail"))
	}

	if productionModeEnabled && conf.PasswordLenMin < 6 {
		return nil, fmt.Errorf(
			"minimum password length must be 6 in production mode, was: %d",
			conf.PasswordLenMin,
		)
	}

	return &validator{
		conf:        conf,
		regexpEmail: regexpEmail,
	}, nil
}
