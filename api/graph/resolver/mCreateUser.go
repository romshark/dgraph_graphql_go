package resolver

import (
	"context"
	"time"

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
) *User {
	// Validate inputs
	if err := rsv.validator.UserDisplayName(params.DisplayName); err != nil {
		err = strerr.Wrap(strerr.ErrInvalidInput, err)
		rsv.error(ctx, err)
		return nil
	}
	if err := rsv.validator.Email(params.Email); err != nil {
		err = strerr.Wrap(strerr.ErrInvalidInput, err)
		rsv.error(ctx, err)
		return nil
	}
	if err := rsv.validator.Password(params.Password); err != nil {
		err = strerr.Wrap(strerr.ErrInvalidInput, err)
		rsv.error(ctx, err)
		return nil
	}

	// Create password hash
	passwordHash, err := rsv.passwordHasher.Hash([]byte(params.Password))
	if err != nil {
		rsv.error(ctx, err)
		return nil
	}

	creationTime := time.Now()

	transactRes, err := rsv.str.CreateUser(
		ctx,
		creationTime,
		params.Email,
		params.DisplayName,
		string(passwordHash),
	)
	if err != nil {
		rsv.error(ctx, err)
		return nil
	}

	return &User{
		root:        rsv,
		uid:         transactRes.UID,
		id:          transactRes.ID,
		creation:    creationTime,
		displayName: params.DisplayName,
		email:       params.Email,
	}
}
