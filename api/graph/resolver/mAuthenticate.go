package resolver

import (
	"context"

	"github.com/romshark/dgraph_graphql_go/api/graph/auth"
	"github.com/romshark/dgraph_graphql_go/store/dgraph"
	strerr "github.com/romshark/dgraph_graphql_go/store/errors"
)

// Authenticate resolves Mutation.authenticate
func (rsv *Resolver) Authenticate(
	ctx context.Context,
	params struct {
		SessionKey string
	},
) *Session {
	var queryResult struct {
		Session []dgraph.Session `json:"session"`
	}
	if err := rsv.str.QueryVars(
		ctx,
		`query Session($sessionKey: string) {
			session(func: eq(Session.key, $sessionKey)) {
				uid
				Session.key
				Session.creation
				Session.user {
					uid
					User.id
				}
			}
		}`,
		map[string]string{"$sessionKey": params.SessionKey},
		&queryResult,
	); err != nil {
		rsv.error(ctx, err)
		return nil
	}

	if len(queryResult.Session) < 1 {
		err := strerr.New(strerr.ErrInvalidInput, "session not found")
		rsv.error(ctx, err)
		return nil
	}

	sess := queryResult.Session[0]

	// Dynamically update the session on successful sign-in
	if session, isSession := ctx.Value(
		auth.CtxSession,
	).(*auth.RequestSession); isSession {
		session.Creation = sess.Creation
		session.UserID = sess.User[0].ID
	}

	return &Session{
		root:     rsv,
		uid:      sess.UID,
		key:      sess.Key,
		creation: sess.Creation,
		userUID:  sess.User[0].UID,
	}
}
