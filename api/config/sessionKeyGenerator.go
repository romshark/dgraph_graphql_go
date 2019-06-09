package config

import (
	"fmt"
	"reflect"
)

// SessionKeyGenerator represents a session key generator identifier
type SessionKeyGenerator string

// UnmarshalTOML implements the TOML unmarshaler interface
func (v *SessionKeyGenerator) UnmarshalTOML(val interface{}) error {
	if str, isString := val.(string); isString {
		switch str {
		case "default":
			*v = SessionKeyGenerator(str)
		default:
			return fmt.Errorf("unknown session key generator: '%s'", val)
		}
		return nil
	}
	return fmt.Errorf(
		"unexpected session key generator value type: %s",
		reflect.TypeOf(val),
	)
}
