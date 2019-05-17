package validator

import (
	"github.com/pkg/errors"
)

// ReactionMessage implements the Validator interface
func (vld *validator) ReactionMessage(v string) error {
	if uint(len(v)) < vld.opts.ReactionMessageLenMin {
		return errors.Errorf(
			"Reaction.message too short (min: %d)",
			vld.opts.ReactionMessageLenMin,
		)
	}
	if uint(len(v)) > vld.opts.ReactionMessageLenMax {
		return errors.Errorf(
			"Reaction.message too long (%d / %d)",
			len(v),
			vld.opts.ReactionMessageLenMax,
		)
	}
	return nil
}
