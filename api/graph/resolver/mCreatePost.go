package resolver

import (
	"context"

	"github.com/pkg/errors"
	"github.com/romshark/dgraph_graphql_go/store"
	"github.com/romshark/dgraph_graphql_go/store/dbmod"
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
	newUID, newID, err := rsv.str.CreatePost(
		ctx,
		store.ID(params.Author),
		params.Title,
		params.Contents,
	)
	if err != nil {
		rsv.error(ctx, err)
		return nil, err
	}

	var result struct {
		NewPost []dbmod.Post `json:"newPost"`
	}
	if err := rsv.str.QueryVars(
		ctx,
		`query NewPost($nodeId: string) {
			newPost(func: uid($nodeId)) {
				Post.creation
				Post.title
				Post.contents
			}
		}`,
		map[string]string{
			"$nodeId": newUID.NodeID,
		},
		&result,
	); err != nil {
		rsv.error(ctx, err)
		return nil, err
	}

	if len(result.NewPost) != 1 {
		err := errors.Errorf(
			"unexpected number of new posts: %d",
			len(result.NewPost),
		)
		rsv.error(ctx, err)
		return nil, err
	}

	newPost := result.NewPost[0]

	return &Post{
		root:     rsv,
		uid:      newUID,
		id:       newID,
		creation: newPost.Creation,
		title:    newPost.Title,
		contents: newPost.Contents,
	}, nil
}
