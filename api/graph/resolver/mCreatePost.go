package resolver

import (
	"context"
	"time"

	"github.com/romshark/dgraph_graphql_go/store"
	"github.com/romshark/dgraph_graphql_go/store/auth"
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
) (*Post, error) {
	if err := auth.Authorize(ctx, auth.IsOwner{
		Owner: store.ID(params.Author),
	}); err != nil {
		rsv.error(ctx, err)
		return nil, err
	}

	// Validate input
	if err := store.ValidatePostTitle(params.Title); err != nil {
		err = strerr.Wrap(strerr.ErrInvalidInput, err)
		rsv.error(ctx, err)
		return nil, err
	}
	if err := store.ValidatePostContents(params.Contents); err != nil {
		err = strerr.Wrap(strerr.ErrInvalidInput, err)
		rsv.error(ctx, err)
		return nil, err
	}

	creationTime := time.Now()

	transactRes, err := rsv.str.CreatePost(
		ctx,
		creationTime,
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
		creation:  creationTime,
		title:     params.Title,
		contents:  params.Contents,
		authorUID: transactRes.AuthorUID,
	}, nil
}
