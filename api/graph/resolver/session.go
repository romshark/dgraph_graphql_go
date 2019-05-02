package resolver

import (
	"context"
	"time"

	"github.com/graph-gophers/graphql-go"
	"github.com/romshark/dgraph_graphql_go/store"
	"github.com/romshark/dgraph_graphql_go/store/dbmod"
)

// Session represents the resolver of the identically named type
type Session struct {
	root     *Resolver
	uid      store.UID
	key      string
	creation time.Time
	userUID  store.UID
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
) (*User, error) {
	var query struct {
		Sessions []dbmod.Session `json:"session"`
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
			"$nodeId": rsv.uid.NodeID,
		},
		&query,
	); err != nil {
		rsv.root.error(ctx, err)
		return nil, err
	}

	owner := query.Sessions[0].User[0]
	return &User{
		root:        rsv.root,
		uid:         store.UID{NodeID: owner.UID},
		id:          owner.ID,
		creation:    owner.Creation,
		email:       owner.Email,
		displayName: owner.DisplayName,
	}, nil
}
