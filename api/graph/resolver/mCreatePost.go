package resolver

import "context"

// CreatePost resolves Mutation.createPost
func (rsv *Resolver) CreatePost(
	ctx context.Context,
	params struct {
		Title    string
		Contents string
	},
) (*Post, error) {
	return nil, nil
}
