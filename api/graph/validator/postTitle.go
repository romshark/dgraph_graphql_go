package validator

import "github.com/pkg/errors"

// PostTitle implements the Valiator interface
func (vld *validator) PostTitle(v string) error {
	if uint(len(v)) < vld.opts.PostTitleLenMin {
		return errors.Errorf(
			"Post.title too short (min: %d)",
			vld.opts.PostTitleLenMin,
		)
	}
	if uint(len(v)) > vld.opts.PostTitleLenMax {
		return errors.Errorf(
			"Post.title too long (%d / %d)",
			len(v),
			vld.opts.PostTitleLenMax,
		)
	}
	return nil
}
