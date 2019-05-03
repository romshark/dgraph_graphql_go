package resolver

import (
	"context"

	"github.com/romshark/dgraph_graphql_go/store"
	"github.com/romshark/dgraph_graphql_go/store/enum/emotion"
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
	emot := emotion.Emotion(params.Emotion)

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
