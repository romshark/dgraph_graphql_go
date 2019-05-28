package resolver

import (
	"context"
	"time"

	"github.com/graph-gophers/graphql-go"
	"github.com/romshark/dgraph_graphql_go/store"
	"github.com/romshark/dgraph_graphql_go/store/dgraph"
)

// Session represents the resolver of the identically named type
type Session struct {
	root     *Resolver
	uid      string
	key      string
	creation time.Time
	userUID  string
}

// Key resolves Session.key
func (rsv *Session) Key() string {
	return rsv.key
}

// Creation resolves Session.creation
func (rsv *Session) Creation() graphql.Time {
	return graphql.Time{Time: rsv.creation}
}

// User resolves Session.user
func (rsv *Session) User(
	ctx context.Context,
) *User {
	var query struct {
		Sessions []dgraph.Session `json:"session"`
	}
	if err := rsv.root.str.QueryVars(
		ctx,
		`query SessionUser($nodeId: string) {
			session(func: uid($nodeId)) {
				Session.user {
					uid
					User.id
					User.creation
					User.email
					User.displayName
				}
			}
		}`,
		map[string]string{
			"$nodeId": rsv.uid,
		},
		&query,
	); err != nil {
		rsv.root.error(ctx, err)
		return nil
	}

	owner := query.Sessions[0].User[0]
	return &User{
		root:        rsv.root,
		uid:         owner.UID,
		id:          store.ID(owner.ID),
		creation:    owner.Creation,
		email:       owner.Email,
		displayName: owner.DisplayName,
	}
}
