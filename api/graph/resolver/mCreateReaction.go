package resolver

import (
	"context"

	"github.com/romshark/dgraph_graphql_go/store"
	"github.com/romshark/dgraph_graphql_go/store/auth"
	"github.com/romshark/dgraph_graphql_go/store/enum/emotion"
	strerr "github.com/romshark/dgraph_graphql_go/store/errors"
)

// CreateReaction resolves Mutation.createReaction
func (rsv *Resolver) CreateReaction(
	ctx context.Context,
	params struct {
		Author  string
		Subject string
		Emotion string
		Message string
	},
) (*Reaction, error) {
	if err := auth.Authorize(ctx, auth.IsOwner{
		Owner: store.ID(params.Author),
	}); err != nil {
		rsv.error(ctx, err)
		return nil, err
	}

	emot := emotion.Emotion(params.Emotion)

	// Validate input
	if err := store.ValidateReactionMessage(params.Message); err != nil {
		err = strerr.Wrap(strerr.ErrInvalidInput, err)
		rsv.error(ctx, err)
		return nil, err
	}
	if err := emotion.Validate(emot); err != nil {
		err = strerr.Wrap(strerr.ErrInvalidInput, err)
		rsv.error(ctx, err)
		return nil, err
	}

	// Create new reaction entity
	transactRes, err := rsv.str.CreateReaction(
		ctx,
		store.ID(params.Author),
		store.ID(params.Subject),
		emot,
		params.Message,
	)
	if err != nil {
		rsv.error(ctx, err)
		return nil, err
	}

	return &Reaction{
		root:       rsv,
		uid:        transactRes.UID,
		id:         transactRes.ID,
		creation:   transactRes.CreationTime,
		authorUID:  transactRes.AuthorUID,
		subjectUID: transactRes.SubjectUID,
		emotion:    emot,
		message:    params.Message,
	}, nil
}
