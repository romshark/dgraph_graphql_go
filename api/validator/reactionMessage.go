package validator

import (
	"github.com/pkg/errors"
)

// ReactionMessage implements the Validator interface
func (vld *validator) ReactionMessage(v string) error {
	if uint(len(v)) < vld.conf.ReactionMessageLenMin {
		return errors.Errorf(
			"Reaction.message too short (min: %d)",
			vld.conf.ReactionMessageLenMin,
		)
	}
	if uint(len(v)) > vld.conf.ReactionMessageLenMax {
		return errors.Errorf(
			"Reaction.message too long (%d / %d)",
			len(v),
			vld.conf.ReactionMessageLenMax,
		)
	}
	return nil
}
