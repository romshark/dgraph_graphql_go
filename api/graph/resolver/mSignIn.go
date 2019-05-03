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
	transactRes, err := rsv.str.CreateSession(
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
		uid:      transactRes.UID,
		key:      transactRes.Key,
		creation: transactRes.CreationTime,
		userUID:  transactRes.UserUID,
	}, nil
}
