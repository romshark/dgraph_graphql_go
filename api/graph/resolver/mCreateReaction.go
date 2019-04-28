package resolver

import "context"

// CreateReaction resolves Mutation.createReaction
func (rsv *Resolver) CreateReaction(
	ctx context.Context,
	params struct {
		Emotion string
		Message string
	},
) (*Reaction, error) {
	return nil, nil
}
