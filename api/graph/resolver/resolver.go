package resolver

import (
	"context"

	"github.com/romshark/dgraph_graphql_go/store"
	"github.com/romshark/dgraph_graphql_go/store/dbmod"
)

// CtxKey represents a context.Context value key type
type CtxKey int

// CtxErrorRef defines the context.Context error reference value key
const CtxErrorRef CtxKey = 1

// Resolver represents the root Graph resolver
type Resolver struct {
	str store.Store
}

// New creates a new graph resolver instance
func New(str store.Store) *Resolver {
	return &Resolver{
		str: str,
	}
}

// Users resolves Query.users
func (rsv *Resolver) Users(ctx context.Context) ([]*User, error) {
	var result struct {
		Users []dbmod.User `json:"users"`
	}
	if err := rsv.str.Query(
		ctx,
		`{
			users(func: has(User.id)) {
				uid
				User.id
				User.creation
				User.email
				User.displayName
			}
		}`,
		&result,
	); err != nil {
		rsv.error(ctx, err)
		return nil, err
	}
	resolvers := make([]*User, len(result.Users))
	for i, usr := range result.Users {
		resolvers[i] = &User{
			root:        rsv,
			uid:         store.UID{NodeID: usr.UID},
			id:          usr.ID,
			displayName: usr.DisplayName,
			email:       usr.Email,
			creation:    usr.Creation,
		}
	}
	return resolvers, nil
}

// Posts resolves Query.posts
func (rsv *Resolver) Posts(ctx context.Context) ([]*Post, error) {
	var result struct {
		Posts []dbmod.Post `json:"posts"`
	}
	if err := rsv.str.Query(
		ctx,
		`{
			posts(func: has(Post.id)) {
				uid
				Post.id
				Post.creation
				Post.title
				Post.contents
			}
		}`,
		&result,
	); err != nil {
		rsv.error(ctx, err)
		return nil, err
	}
	resolvers := make([]*Post, len(result.Posts))
	for i, usr := range result.Posts {
		resolvers[i] = &Post{
			root:     rsv,
			uid:      store.UID{NodeID: usr.UID},
			id:       usr.ID,
			title:    usr.Title,
			contents: usr.Contents,
			creation: usr.Creation,
		}
	}
	return resolvers, nil
}

// User resolves Query.user
func (rsv *Resolver) User(
	ctx context.Context,
	param struct {
		Id string
	},
) (*User, error) {
	var result struct {
		Users []dbmod.User `json:"users"`
	}
	if err := rsv.str.QueryVars(
		ctx,
		`query User($userId: string) {
			users(func: eq(User.id, $userId)) {
				uid
				User.id
				User.creation
				User.email
				User.displayName
			}
		}`,
		map[string]string{
			"$userId": param.Id,
		},
		&result,
	); err != nil {
		rsv.error(ctx, err)
		return nil, err
	}
	if len(result.Users) < 1 {
		return nil, nil
	}

	usr := result.Users[0]
	return &User{
		root:        rsv,
		uid:         store.UID{NodeID: usr.UID},
		id:          usr.ID,
		displayName: usr.DisplayName,
		email:       usr.Email,
		creation:    usr.Creation,
	}, nil
}

// Post resolves Query.post
func (rsv *Resolver) Post(
	ctx context.Context,
	param struct {
		Id string
	},
) (*Post, error) {
	var result struct {
		Posts []dbmod.Post `json:"posts"`
	}
	if err := rsv.str.QueryVars(
		ctx,
		`query Post($postId: string) {
			posts(func: eq(Post.id, $postId)) {
				uid
				Post.id
				Post.creation
				Post.title
				Post.contents
			}
		}`,
		map[string]string{
			"$postId": param.Id,
		},
		&result,
	); err != nil {
		rsv.error(ctx, err)
		return nil, err
	}
	if len(result.Posts) < 1 {
		return nil, nil
	}

	usr := result.Posts[0]
	return &Post{
		root:     rsv,
		uid:      store.UID{NodeID: usr.UID},
		id:       usr.ID,
		title:    usr.Title,
		contents: usr.Contents,
		creation: usr.Creation,
	}, nil
}

// error writes an error to the resolver context for the API server to read
func (rsv *Resolver) error(ctx context.Context, err error) {
	ctxErr := ctx.Value(CtxErrorRef).(*error)
	*ctxErr = err
}
