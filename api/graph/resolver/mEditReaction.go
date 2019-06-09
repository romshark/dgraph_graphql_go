package resolver

import (
	"context"

	"github.com/romshark/dgraph_graphql_go/api/graph/auth"
	"github.com/romshark/dgraph_graphql_go/store"
	strerr "github.com/romshark/dgraph_graphql_go/store/errors"
)

// EditReaction resolves Mutation.editReaction
func (rsv *Resolver) EditReaction(
	ctx context.Context,
	params struct {
		Reaction   string
		Editor     string
		NewMessage string
	},
) *Reaction {
	if err := auth.Authorize(ctx, auth.IsOwner{
		Owner: store.ID(params.Editor),
	}); err != nil {
		rsv.error(ctx, err)
		return nil
	}

	// Validate input
	if err := rsv.validator.ReactionMessage(params.NewMessage); err != nil {
		err = strerr.Wrap(strerr.ErrInvalidInput, err)
		rsv.error(ctx, err)
		return nil
	}

	mutatedReaction, _, err := rsv.str.EditReaction(
		ctx,
		store.ID(params.Reaction),
		store.ID(params.Editor),
		params.NewMessage,
	)
	if err != nil {
		rsv.error(ctx, err)
		return nil
	}

	return &Reaction{
		root:       rsv,
		uid:        mutatedReaction.UID,
		id:         store.ID(params.Reaction),
		creation:   mutatedReaction.Creation,
		emotion:    mutatedReaction.Emotion,
		message:    mutatedReaction.Message,
		authorUID:  mutatedReaction.Author.UID,
		subjectUID: mutatedReaction.Subject.NodeID(),
	}
}
