package dgraph

import (
	"context"
	"encoding/json"

	"github.com/dgraph-io/dgo/protos/api"
	"github.com/romshark/dgraph_graphql_go/store"
	strerr "github.com/romshark/dgraph_graphql_go/store/errors"
)

// CloseAllSessions closes all sessions of the given user
func (str *impl) CloseAllSessions(
	ctx context.Context,
	user store.ID,
) (
	result []string,
	err error,
) {
	// Begin transaction
	txn, close := str.txn(&err)
	if err != nil {
		return
	}
	defer close()

	// Find the user and all associated sessions
	var qr struct {
		User []User `json:"user"`
	}
	err = txn.QueryVars(
		ctx,
		`query Sessions(
			$userID: string
		) {
			user(func: eq(User.id, $userID)) {
				uid
				User.sessions {
					uid
					Session.key
					~sessions {
						uid
					}
				}
			}
		}`,
		map[string]string{
			"$userID": string(user),
		},
		&qr,
	)
	if err != nil {
		return
	}

	if len(qr.User) < 1 {
		err = strerr.Newf(
			strerr.ErrInvalidInput,
			"user not found",
		)
		return
	}

	usr := qr.User[0]
	deletions := make([]interface{}, 0, len(usr.Sessions)*2+1)

	result = make([]string, len(usr.Sessions))
	for i, sess := range usr.Sessions {
		result[i] = sess.Key

		// Delete the global "sessions" references
		deletions = append(deletions, sess.RSessions[0])

		// Delete the actual Session nodes
		deletions = append(deletions, UID{NodeID: sess.UID})
	}

	// Delete the "User.sessions" references
	userSessions := make([]UID, len(usr.Sessions))
	for i, sess := range usr.Sessions {
		userSessions[i] = UID{NodeID: sess.UID}
	}
	deletions = append(deletions, struct {
		UID          string `json:"uid"`
		UserSessions []UID  `json:"User.sessions"`
	}{
		UID:          usr.UID,
		UserSessions: userSessions,
	})

	var deleteJSON []byte
	deleteJSON, err = json.Marshal(deletions)
	if err != nil {
		return
	}
	_, err = txn.Mutation(ctx, &api.Mutation{DeleteJson: deleteJSON})
	return
}
