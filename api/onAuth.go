package api

import (
	"context"
	"time"

	"github.com/romshark/dgraph_graphql_go/store"
	"github.com/romshark/dgraph_graphql_go/store/dgraph"
)

// onAuth is invoked by the transport layer during client authentication
func (srv *server) onAuth(
	ctx context.Context,
	sessionKey string,
) (userID store.ID, sessionCreationTime time.Time) {
	// Search for the user session by key
	var result struct {
		Session []dgraph.Session `json:"session"`
	}
	if err := srv.store.QueryVars(
		ctx,
		`query Session($sessionKey: string) {
			session(func: eq(Session.key, $sessionKey)) {
				Session.creation
				Session.user {
					uid
					User.id
				}
			}
		}`,
		map[string]string{
			"$sessionKey": sessionKey,
		},
		&result,
	); err != nil {
		return
	}

	if len(result.Session) < 1 {
		return
	}

	userID = store.ID(result.Session[0].User[0].ID)
	sessionCreationTime = result.Session[0].Creation
	return
}
