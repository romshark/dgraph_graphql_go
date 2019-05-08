package resolver

import (
	"context"

	"github.com/romshark/dgraph_graphql_go/store"
	strerr "github.com/romshark/dgraph_graphql_go/store/errors"
)

// CreateUser resolves Mutation.createUser
func (rsv *Resolver) CreateUser(
	ctx context.Context,
	params struct {
		Email       string
		DisplayName string
		Password    string
	},
) (*User, error) {
	// Validate inputs
	if err := store.ValidateUserDisplayName(params.DisplayName); err != nil {
		err = strerr.Wrap(strerr.ErrInvalidInput, err)
		rsv.error(ctx, err)
		return nil, err
	}
	if err := store.ValidateEmail(params.Email); err != nil {
		err = strerr.Wrap(strerr.ErrInvalidInput, err)
		rsv.error(ctx, err)
		return nil, err
	}
	if err := store.ValidatePassword(params.Password); err != nil {
		err = strerr.Wrap(strerr.ErrInvalidInput, err)
		rsv.error(ctx, err)
		return nil, err
	}

	transactRes, err := rsv.str.CreateUser(
		ctx,
		params.Email,
		params.DisplayName,
		params.Password,
	)
	if err != nil {
		rsv.error(ctx, err)
		return nil, err
	}

	return &User{
		root:        rsv,
		uid:         transactRes.UID,
		id:          transactRes.ID,
		creation:    transactRes.CreationTime,
		displayName: params.DisplayName,
		email:       params.Email,
	}, nil
}
