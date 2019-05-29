package validator

import "github.com/pkg/errors"

// PostContents implements the Validator interface
func (vld *validator) PostContents(v string) error {
	if uint(len(v)) < vld.conf.PostContentsLenMin {
		return errors.Errorf(
			"Post.contents too short (min: %d)",
			vld.conf.PostContentsLenMin,
		)
	}
	if uint(len(v)) > vld.conf.PostContentsLenMax {
		return errors.Errorf(
			"Post.contents too long (%d / %d)",
			len(v),
			vld.conf.PostContentsLenMax,
		)
	}
	return nil
}
