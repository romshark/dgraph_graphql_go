package validator

import "github.com/pkg/errors"

// PostContents implements the Validator interface
func (vld *validator) PostContents(v string) error {
	if uint(len(v)) < vld.opts.PostContentsLenMin {
		return errors.Errorf(
			"Post.contents too short (min: %d)",
			vld.opts.PostContentsLenMin,
		)
	}
	if uint(len(v)) > vld.opts.PostContentsLenMax {
		return errors.Errorf(
			"Post.contents too long (%d / %d)",
			len(v),
			vld.opts.PostContentsLenMax,
		)
	}
	return nil
}
