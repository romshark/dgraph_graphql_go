package store

import "github.com/pkg/errors"

// ValidatePostContents returns an error if invalid, otherwise returns nil
func ValidatePostContents(v string) error {
	if len(v) < 1 {
		return errors.Errorf("Post.contents too short (min: 1)")
	}
	if len(v) > 256 {
		return errors.Errorf("Post.contents too long (%d / 256)", len(v))
	}
	return nil
}
