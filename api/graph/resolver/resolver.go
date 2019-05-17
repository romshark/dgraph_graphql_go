package resolver

import (
	"context"

	"github.com/pkg/errors"
	"github.com/romshark/dgraph_graphql_go/api/passhash"
	"github.com/romshark/dgraph_graphql_go/api/sesskeygen"
	"github.com/romshark/dgraph_graphql_go/api/validator"
	"github.com/romshark/dgraph_graphql_go/store"
	"github.com/romshark/dgraph_graphql_go/store/dgraph"
)

// CtxKey represents a context.Context value key type
type CtxKey int

// CtxErrorRef defines the context.Context error reference value key
const CtxErrorRef CtxKey = 1

// Resolver represents the root Graph resolver
type Resolver struct {
	str                 store.Store
	validator           validator.Validator
	sessionKeyGenerator sesskeygen.SessionKeyGenerator
	passwordHasher      passhash.PasswordHasher
}

// New creates a new graph resolver instance
func New(
	str store.Store,
	validator validator.Validator,
	sessionKeyGenerator sesskeygen.SessionKeyGenerator,
	passwordHasher passhash.PasswordHasher,
) (*Resolver, error) {
	if sessionKeyGenerator == nil {
		return nil, errors.Errorf(
			"missing session key generator during resolver initialization",
		)
	}
	if passwordHasher == nil {
		return nil, errors.Errorf(
			"missing password hasher during resolver initialization",
		)
	}

	return &Resolver{
		str:                 str,
		validator:           validator,
		sessionKeyGenerator: sessionKeyGenerator,
		passwordHasher:      passwordHasher,
	}, nil
}

// Users resolves Query.users
func (rsv *Resolver) Users(ctx context.Context) ([]*User, error) {
	var result struct {
		Users []dgraph.User `json:"users"`
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
			uid:         usr.UID,
			id:          store.ID(usr.ID),
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
		Posts []dgraph.Post `json:"posts"`
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
				Post.author {
					uid
				}
			}
		}`,
		&result,
	); err != nil {
		rsv.error(ctx, err)
		return nil, err
	}
	resolvers := make([]*Post, len(result.Posts))
	for i, post := range result.Posts {
		resolvers[i] = &Post{
			root:      rsv,
			uid:       post.UID,
			id:        post.ID,
			title:     post.Title,
			contents:  post.Contents,
			creation:  post.Creation,
			authorUID: post.UID,
		}
	}
	return resolvers, nil
}

// User resolves Query.user
func (rsv *Resolver) User(
	ctx context.Context,
	params struct {
		ID string
	},
) (*User, error) {
	var result struct {
		Users []dgraph.User `json:"users"`
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
			"$userId": params.ID,
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
		uid:         usr.UID,
		id:          store.ID(usr.ID),
		displayName: usr.DisplayName,
		email:       usr.Email,
		creation:    usr.Creation,
	}, nil
}

// Post resolves Query.post
func (rsv *Resolver) Post(
	ctx context.Context,
	params struct {
		ID string
	},
) (*Post, error) {
	var result struct {
		Posts []dgraph.Post `json:"posts"`
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
				Post.author {
					uid
				}
			}
		}`,
		map[string]string{
			"$postId": params.ID,
		},
		&result,
	); err != nil {
		rsv.error(ctx, err)
		return nil, err
	}
	if len(result.Posts) < 1 {
		return nil, nil
	}

	post := result.Posts[0]
	return &Post{
		root:      rsv,
		uid:       post.UID,
		id:        post.ID,
		title:     post.Title,
		contents:  post.Contents,
		creation:  post.Creation,
		authorUID: post.Author[0].UID,
	}, nil
}

// Reaction resolves Query.reaction
func (rsv *Resolver) Reaction(
	ctx context.Context,
	params struct {
		ID string
	},
) (*Reaction, error) {
	var result struct {
		Reactions []dgraph.Reaction `json:"reactions"`
	}
	if err := rsv.str.QueryVars(
		ctx,
		`query Reaction($reactionId: string) {
			reactions(func: eq(Reaction.id, $reactionId)) {
				uid
				Reaction.id
				Reaction.creation
				Reaction.emotion
				Reaction.message
				Reaction.author {
					uid
				}
				Reaction.subject {
					uid
					Post.id
					Reaction.id
				}
			}
		}`,
		map[string]string{
			"$reactionId": params.ID,
		},
		&result,
	); err != nil {
		rsv.error(ctx, err)
		return nil, err
	}
	if len(result.Reactions) < 1 {
		return nil, nil
	}

	reaction := result.Reactions[0]
	return &Reaction{
		root:       rsv,
		uid:        reaction.UID,
		id:         reaction.ID,
		emotion:    reaction.Emotion,
		message:    reaction.Message,
		creation:   reaction.Creation,
		authorUID:  reaction.Author[0].UID,
		subjectUID: *reaction.Subject[0].UID(),
	}, nil
}

// error writes an error to the resolver context for the API server to read
func (rsv *Resolver) error(ctx context.Context, err error) {
	ctxErr := ctx.Value(CtxErrorRef).(*error)
	*ctxErr = err
}
