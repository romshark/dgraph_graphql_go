package validator

import "github.com/pkg/errors"

// PostTitle implements the Valiator interface
func (vld *validator) PostTitle(v string) error {
	if uint(len(v)) < vld.conf.PostTitleLenMin {
		return errors.Errorf(
			"Post.title too short (min: %d)",
			vld.conf.PostTitleLenMin,
		)
	}
	if uint(len(v)) > vld.conf.PostTitleLenMax {
		return errors.Errorf(
			"Post.title too long (%d / %d)",
			len(v),
			vld.conf.PostTitleLenMax,
		)
	}
	return nil
}
