package gqlshield

import "github.com/pkg/errors"

func validateParameterName(name string) error {
	if len(name) < 1 {
		return errors.New("invalid query parameter name (empty)")
	}
	return nil
}
