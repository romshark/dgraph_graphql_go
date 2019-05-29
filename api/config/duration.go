package config

import (
	"fmt"
	"reflect"
	"time"
)

// Duration represents a duration
type Duration time.Duration

// UnmarshalTOML implements the TOML unmarshaler interface
func (v *Duration) UnmarshalTOML(val interface{}) error {
	if str, isString := val.(string); isString {
		dur, err := time.ParseDuration(str)
		if err != nil {
			return err
		}
		*v = Duration(dur)
		return nil
	}
	return fmt.Errorf(
		"unexpected duration value type: %s",
		reflect.TypeOf(val),
	)
}
