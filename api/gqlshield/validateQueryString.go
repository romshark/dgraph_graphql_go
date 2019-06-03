package gqlshield

import "github.com/pkg/errors"

func validateQueryString(queryString string) error {
	if len(queryString) < 1 {
		return errors.New("invalid query string (empty)")
	}
	return nil
}
