package store

import "github.com/pkg/errors"

// ValidatePostTitle returns an error if invalid, otherwise returns nil
func ValidatePostTitle(v string) error {
	if len(v) < 2 {
		return errors.Errorf("Post.title too short (min: 2)")
	}
	if len(v) > 64 {
		return errors.Errorf("Post.title too long (%d / 64)", len(v))
	}
	return nil
}
