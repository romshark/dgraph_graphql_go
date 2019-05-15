package auth

import (
	"context"
	"time"

	"github.com/romshark/dgraph_graphql_go/store"
	"github.com/romshark/dgraph_graphql_go/store/errors"
)

// CtxKey represents a context.Context value key type
type CtxKey int

// CtxSession defines the context.Context session value key
const CtxSession CtxKey = 1

// RequestSession represents a client session
type RequestSession struct {
	IsDebug  bool
	UserID   store.ID
	Creation time.Time
}

// Requirement defines the authorization requirement implementation interface
type Requirement interface {
	check(session *RequestSession) string
}

// Authorize authorizes the client session taken from the provided context
// against the provided requirements
func Authorize(
	ctx context.Context,
	requirements ...Requirement,
) error {
	// Extract the session
	session, isSession := ctx.Value(CtxSession).(*RequestSession)

	// Pass debug clients without further authorization
	if isSession && session.IsDebug {
		return nil
	}

	// Check all requirements
	for _, rule := range requirements {
		if errMsg := rule.check(session); len(errMsg) > 0 {
			// Unauthorized
			return errors.New(errors.ErrUnauthorized, errMsg)
		}
	}

	// Authorized
	return nil
}
