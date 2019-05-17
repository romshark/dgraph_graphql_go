package resolver

import (
	"context"
	"time"

	"github.com/romshark/dgraph_graphql_go/api/graph/auth"
	strerr "github.com/romshark/dgraph_graphql_go/store/errors"
)

// SignIn resolves Mutation.signIn
func (rsv *Resolver) SignIn(
	ctx context.Context,
	params struct {
		Email    string
		Password string
	},
) (*Session, error) {
	// Validate inputs
	if len(params.Email) < 1 || len(params.Password) < 1 {
		err := strerr.New(strerr.ErrInvalidInput, "missing credentials")
		rsv.error(ctx, err)
		return nil, err
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
		return nil, err
	}

	// Dynamically update the session on successful sign-in
	if session, isSession := ctx.Value(
		auth.CtxSession,
	).(*auth.RequestSession); isSession {
		session.Creation = creationTime
		session.UserID = newSession.User.ID
	}

	return &Session{
		root:     rsv,
		uid:      newSession.UID,
		key:      key,
		creation: creationTime,
		userUID:  newSession.User.UID,
	}, nil
}
