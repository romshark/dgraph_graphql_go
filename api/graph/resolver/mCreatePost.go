package resolver

import (
	"context"

	"github.com/romshark/dgraph_graphql_go/store"
)

// CreatePost resolves Mutation.createPost
func (rsv *Resolver) CreatePost(
	ctx context.Context,
	params struct {
		Author   string
		Title    string
		Contents string
	},
) (*Post, error) {
	transactRes, err := rsv.str.CreatePost(
		ctx,
		store.ID(params.Author),
		params.Title,
		params.Contents,
	)
	if err != nil {
		rsv.error(ctx, err)
		return nil, err
	}

	return &Post{
		root:      rsv,
		uid:       transactRes.UID,
		id:        transactRes.ID,
		creation:  transactRes.CreationTime,
		title:     params.Title,
		contents:  params.Contents,
		authorUID: transactRes.AuthorUID,
	}, nil
}
