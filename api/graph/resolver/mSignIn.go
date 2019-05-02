package resolver

import (
	"context"
)

// SignIn resolves Mutation.signIn
func (rsv *Resolver) SignIn(
	ctx context.Context,
	params struct {
		Email    string
		Password string
	},
) (*Session, error) {
	newUID, key, creation, userUID, err := rsv.str.CreateSession(
		ctx,
		params.Email,
		params.Password,
	)
	if err != nil {
		rsv.error(ctx, err)
		return nil, err
	}

	return &Session{
		root:     rsv,
		uid:      newUID,
		key:      key,
		creation: creation,
		userUID:  userUID,
	}, nil
}
