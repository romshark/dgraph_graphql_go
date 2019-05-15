package resolver

import (
	"context"
	"time"

	"github.com/graph-gophers/graphql-go"
	"github.com/romshark/dgraph_graphql_go/store"
	"github.com/romshark/dgraph_graphql_go/store/auth"
	"github.com/romshark/dgraph_graphql_go/store/dgraph"
)

// User represents the resolver of the identically named type
type User struct {
	root        *Resolver
	uid         string
	id          store.ID
	creation    time.Time
	email       string
	displayName string
}

// ID resolves User.id
func (rsv *User) ID() store.ID {
	return rsv.id
}

// Creation resolves User.creation
func (rsv *User) Creation() graphql.Time {
	return graphql.Time{
		Time: rsv.creation,
	}
}

// Email resolves User.email
func (rsv *User) Email(ctx context.Context) (string, error) {
	// Check permissions
	if err := auth.Authorize(ctx, auth.IsOwner{
		Owner: store.ID(rsv.id),
	}); err != nil {
		rsv.root.error(ctx, err)
		return "", err
	}

	return rsv.email, nil
}

// DisplayName resolves User.displayName
func (rsv *User) DisplayName() string {
	return rsv.displayName
}

// Posts resolves User.posts
func (rsv *User) Posts(
	ctx context.Context,
) ([]*Post, error) {
	var query struct {
		Users []dgraph.User `json:"users"`
	}
	if err := rsv.root.str.QueryVars(
		ctx,
		`query Posts($nodeId: string) {
			users(func: uid($nodeId)) {
				User.posts {
					uid
					Post.id
					Post.creation
					Post.title
					Post.contents
				}
			}
		}`,
		map[string]string{
			"$nodeId": rsv.uid,
		},
		&query,
	); err != nil {
		rsv.root.error(ctx, err)
		return nil, err
	}

	if len(query.Users) < 1 {
		return nil, nil
	}

	usr := query.Users[0]
	resolvers := make([]*Post, len(usr.Posts))
	for i, post := range usr.Posts {
		resolvers[i] = &Post{
			root:      rsv.root,
			uid:       post.UID,
			id:        post.ID,
			creation:  post.Creation,
			title:     post.Title,
			contents:  post.Contents,
			authorUID: rsv.uid,
		}
	}

	return resolvers, nil
}

// Sessions resolves User.sessions
func (rsv *User) Sessions(
	ctx context.Context,
) ([]*Session, error) {
	// Check permissions
	if err := auth.Authorize(ctx, auth.IsOwner{
		Owner: store.ID(rsv.id),
	}); err != nil {
		rsv.root.error(ctx, err)
		return nil, err
	}

	var query struct {
		Users []dgraph.User `json:"users"`
	}
	if err := rsv.root.str.QueryVars(
		ctx,
		`query Sessions($nodeId: string) {
			users(func: uid($nodeId)) {
				User.sessions {
					uid
					Session.key
					Session.creation
				}
			}
		}`,
		map[string]string{
			"$nodeId": rsv.uid,
		},
		&query,
	); err != nil {
		rsv.root.error(ctx, err)
		return nil, err
	}

	if len(query.Users) < 1 {
		return nil, nil
	}

	usr := query.Users[0]
	resolvers := make([]*Session, len(usr.Sessions))
	for i, sess := range usr.Sessions {
		resolvers[i] = &Session{
			root:     rsv.root,
			uid:      sess.UID,
			key:      sess.Key,
			creation: sess.Creation,
			userUID:  rsv.uid,
		}
	}

	return resolvers, nil
}

// PublishedReactions resolves User.publishedReactions
func (rsv *User) PublishedReactions(
	ctx context.Context,
) ([]*Reaction, error) {
	var query struct {
		Users []dgraph.User `json:"users"`
	}
	if err := rsv.root.str.QueryVars(
		ctx,
		`query PublishedReactions($nodeId: string) {
			users(func: uid($nodeId)) {
				User.publishedReactions {
					uid
					Reaction.id
					Reaction.creation
					Reaction.emotion
					Reaction.message
					Reaction.subject {
						uid
						Post.id
						Reaction.id
					}
				}
			}
		}`,
		map[string]string{
			"$nodeId": rsv.uid,
		},
		&query,
	); err != nil {
		rsv.root.error(ctx, err)
		return nil, err
	}

	if len(query.Users) < 1 {
		return nil, nil
	}

	usr := query.Users[0]
	resolvers := make([]*Reaction, len(usr.PublishedReactions))
	for i, reaction := range usr.PublishedReactions {
		var subjectUID string
		switch v := reaction.Subject[0].V.(type) {
		case *dgraph.Post:
			subjectUID = v.UID
		case *dgraph.Reaction:
			subjectUID = v.UID
		}
		resolvers[i] = &Reaction{
			root:       rsv.root,
			uid:        reaction.UID,
			id:         reaction.ID,
			subjectUID: subjectUID,
			authorUID:  rsv.uid,
			creation:   reaction.Creation,
			emotion:    reaction.Emotion,
			message:    reaction.Message,
		}
	}

	return resolvers, nil
}
