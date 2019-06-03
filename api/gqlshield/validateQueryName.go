package gqlshield

import "github.com/pkg/errors"

func validateQueryName(name string) error {
	if len(name) < 1 {
		return errors.New("invalid query name (empty)")
	}
	return nil
}
