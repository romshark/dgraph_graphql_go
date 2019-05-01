package resolver

import (
	"context"
	"time"

	"github.com/graph-gophers/graphql-go"
	"github.com/romshark/dgraph_graphql_go/store"
	"github.com/romshark/dgraph_graphql_go/store/dbmod"
)

// User represents the resolver of the identically named type
type User struct {
	root        *Resolver
	uid         store.UID
	id          store.ID
	creation    time.Time
	email       string
	displayName string
}

// Id resolves User.id
func (rsv *User) Id() store.ID {
	return rsv.id
}

// Creation resolves User.creation
func (rsv *User) Creation() graphql.Time {
	return graphql.Time{
		Time: rsv.creation,
	}
}

// Email resolves User.email
func (rsv *User) Email() string {
	return rsv.email
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
		Users []dbmod.User `json:"users"`
	}
	if err := rsv.root.str.QueryVars(
		ctx,
		`query Posts($nodeId: string) {
			users(func: uid($nodeId)) {
				User.posts {
					uid
					Post.id
					Post.creation
					Post.author
					Post.title
					Post.contents
					Post.reaction
				}
			}
		}`,
		map[string]string{
			"$nodeId": rsv.uid.NodeID,
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
			root:     rsv.root,
			uid:      store.UID{NodeID: post.UID},
			id:       post.ID,
			creation: post.Creation,
			title:    post.Title,
			contents: post.Contents,
		}
	}

	return resolvers, nil
}
