package dgraph

import (
	"context"
	"encoding/json"

	"github.com/dgraph-io/dgo/v2/protos/api"
	"github.com/romshark/dgraph_graphql_go/api/graph/auth"
	strerr "github.com/romshark/dgraph_graphql_go/store/errors"
)

// CloseSession closes the given session
func (str *impl) CloseSession(
	ctx context.Context,
	key string,
) (
	result bool,
	err error,
) {
	result = true

	// Begin transaction
	txn, close := str.txn(&err)
	if err != nil {
		return
	}
	defer close()

	// Find the session, its owner and the global reference node
	var qr struct {
		Session []Session `json:"session"`
	}
	err = txn.QueryVars(
		ctx,
		`query Session(
			$key: string
		) {
			session(func: eq(Session.key, $key)) {
				uid
				Session.user {
					uid
					User.id
				}
				~sessions {
					uid
				}
			}
		}`,
		map[string]string{
			"$key": key,
		},
		&qr,
	)
	if err != nil {
		return
	}

	if len(qr.Session) < 1 {
		err = strerr.Newf(
			strerr.ErrInvalidInput,
			"session not found",
		)
		return
	}

	sess := qr.Session[0]

	// Authorize client
	if err = auth.Authorize(ctx, auth.IsOwner{
		Owner: sess.User[0].ID,
	}); err != nil {
		return
	}

	var deleteJSON []byte
	deleteJSON, err = json.Marshal([]interface{}{
		// Delete the global "sessions" reference
		sess.RSessions[0],

		// Delete the "User.sessions" reference
		struct {
			UID          string `json:"uid"`
			UserSessions []UID  `json:"User.sessions"`
		}{
			UID:          sess.User[0].UID,
			UserSessions: []UID{UID{NodeID: sess.UID}},
		},

		// Delete the actual Session node
		UID{NodeID: sess.UID},
	})
	if err != nil {
		return
	}
	_, err = txn.Mutation(ctx, &api.Mutation{DeleteJson: deleteJSON})
	return
}
