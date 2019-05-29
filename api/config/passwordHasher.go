package config

import (
	"fmt"
	"reflect"
)

// PasswordHasher represents a password hasher identifier
type PasswordHasher string

// UnmarshalTOML implements the TOML unmarshaler interface
func (v *PasswordHasher) UnmarshalTOML(val interface{}) error {
	if str, isString := val.(string); isString {
		switch str {
		case "bcrypt":
			*v = PasswordHasher(str)
		default:
			return fmt.Errorf("unknown password hasher: '%s'", val)
		}
		return nil
	}
	return fmt.Errorf(
		"unexpected password hasher value type: %s",
		reflect.TypeOf(val),
	)
}
