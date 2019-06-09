package resolver

import (
	"context"
	"time"

	"github.com/romshark/dgraph_graphql_go/api/graph/auth"
	"github.com/romshark/dgraph_graphql_go/store"
	strerr "github.com/romshark/dgraph_graphql_go/store/errors"
)

// CreatePost resolves Mutation.createPost
func (rsv *Resolver) CreatePost(
	ctx context.Context,
	params struct {
		Author   string
		Title    string
		Contents string
	},
) *Post {
	if err := auth.Authorize(ctx, auth.IsOwner{
		Owner: store.ID(params.Author),
	}); err != nil {
		rsv.error(ctx, err)
		return nil
	}

	// Validate input
	if err := rsv.validator.PostTitle(params.Title); err != nil {
		err = strerr.Wrap(strerr.ErrInvalidInput, err)
		rsv.error(ctx, err)
		return nil
	}
	if err := rsv.validator.PostContents(params.Contents); err != nil {
		err = strerr.Wrap(strerr.ErrInvalidInput, err)
		rsv.error(ctx, err)
		return nil
	}

	creationTime := time.Now()

	newPost, err := rsv.str.CreatePost(
		ctx,
		creationTime,
		store.ID(params.Author),
		params.Title,
		params.Contents,
	)
	if err != nil {
		rsv.error(ctx, err)
		return nil
	}

	return &Post{
		root:      rsv,
		uid:       newPost.UID,
		id:        newPost.ID,
		creation:  creationTime,
		title:     params.Title,
		contents:  params.Contents,
		authorUID: newPost.Author.UID,
	}
}
