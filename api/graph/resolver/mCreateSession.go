package resolver

import (
	"context"
	"time"

	"github.com/romshark/dgraph_graphql_go/api/graph/auth"
	strerr "github.com/romshark/dgraph_graphql_go/store/errors"
)

// CreateSession resolves Mutation.createSession
func (rsv *Resolver) CreateSession(
	ctx context.Context,
	params struct {
		Email    string
		Password string
	},
) *Session {
	// Validate inputs
	if len(params.Email) < 1 || len(params.Password) < 1 {
		err := strerr.New(strerr.ErrInvalidInput, "missing credentials")
		rsv.error(ctx, err)
		return nil
	}

	// Generate session key
	key := rsv.sessionKeyGenerator.Generate()
	creationTime := time.Now()

	newSession, err := rsv.str.CreateSession(
		ctx,
		key,
		creationTime,
		params.Email,
		params.Password,
	)
	if err != nil {
		rsv.error(ctx, err)
		return nil
	}

	// Dynamically update the session on successful sign-in
	if session, isSession := ctx.Value(
		auth.CtxSession,
	).(*auth.RequestSession); isSession {
		session.Creation = creationTime
		session.UserID = newSession.User.ID
		session.ShieldClientRole = auth.GQLShieldClientRegular
	}

	return &Session{
		root:     rsv,
		uid:      newSession.UID,
		key:      key,
		creation: creationTime,
		userUID:  newSession.User.UID,
	}
}
