package resolver

import (
	"context"

	"github.com/romshark/dgraph_graphql_go/api/graph/auth"
	"github.com/romshark/dgraph_graphql_go/store"
	strerr "github.com/romshark/dgraph_graphql_go/store/errors"
)

// EditPost resolves Mutation.editPost
func (rsv *Resolver) EditPost(
	ctx context.Context,
	params struct {
		Post        string
		Editor      string
		NewTitle    *string
		NewContents *string
	},
) *Post {
	if err := auth.Authorize(ctx, auth.IsOwner{
		Owner: store.ID(params.Editor),
	}); err != nil {
		rsv.error(ctx, err)
		return nil
	}

	// Validate input
	if params.NewTitle == nil && params.NewContents == nil {
		err := strerr.New(strerr.ErrInvalidInput, "no changes")
		rsv.error(ctx, err)
		return nil
	}
	if params.NewTitle != nil {
		if err := rsv.validator.PostTitle(*params.NewTitle); err != nil {
			err = strerr.Wrap(strerr.ErrInvalidInput, err)
			rsv.error(ctx, err)
			return nil
		}
	}
	if params.NewContents != nil {
		if err := rsv.validator.PostContents(*params.NewContents); err != nil {
			err = strerr.Wrap(strerr.ErrInvalidInput, err)
			rsv.error(ctx, err)
			return nil
		}
	}

	mutatedPost, _, err := rsv.str.EditPost(
		ctx,
		store.ID(params.Post),
		store.ID(params.Editor),
		params.NewTitle,
		params.NewContents,
	)
	if err != nil {
		rsv.error(ctx, err)
		return nil
	}

	return &Post{
		root:      rsv,
		uid:       mutatedPost.UID,
		id:        store.ID(params.Post),
		creation:  mutatedPost.Creation,
		title:     mutatedPost.Title,
		contents:  mutatedPost.Contents,
		authorUID: mutatedPost.Author.UID,
	}
}
