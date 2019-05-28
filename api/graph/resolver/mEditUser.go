package resolver

import (
	"context"

	"github.com/romshark/dgraph_graphql_go/api/graph/auth"
	"github.com/romshark/dgraph_graphql_go/store"
	strerr "github.com/romshark/dgraph_graphql_go/store/errors"
)

// EditUser resolves Mutation.editUser
func (rsv *Resolver) EditUser(
	ctx context.Context,
	params struct {
		User        string
		Editor      string
		NewEmail    *string
		NewPassword *string
	},
) *User {
	if err := auth.Authorize(ctx, auth.IsOwner{
		Owner: store.ID(params.Editor),
	}); err != nil {
		rsv.error(ctx, err)
		return nil
	}
	if err := auth.Authorize(ctx, auth.IsOwner{
		Owner: store.ID(params.User),
	}); err != nil {
		rsv.error(ctx, err)
		return nil
	}

	// Validate input
	if params.NewEmail == nil && params.NewPassword == nil {
		err := strerr.New(strerr.ErrInvalidInput, "no changes")
		rsv.error(ctx, err)
		return nil
	}
	if params.NewEmail != nil {
		if err := rsv.validator.Email(*params.NewEmail); err != nil {
			err = strerr.Wrap(strerr.ErrInvalidInput, err)
			rsv.error(ctx, err)
			return nil
		}
	}
	if params.NewPassword != nil {
		if err := rsv.validator.Password(*params.NewPassword); err != nil {
			err = strerr.Wrap(strerr.ErrInvalidInput, err)
			rsv.error(ctx, err)
			return nil
		}
	}

	// Create password hash if any
	if params.NewPassword != nil {
		passwordHash, err := rsv.passwordHasher.Hash(
			[]byte(*params.NewPassword),
		)
		if err != nil {
			rsv.error(ctx, err)
			return nil
		}
		*params.NewPassword = string(passwordHash)
	}

	mutatedUser, _, err := rsv.str.EditUser(
		ctx,
		store.ID(params.User),
		store.ID(params.Editor),
		params.NewEmail,
		params.NewPassword,
	)
	if err != nil {
		rsv.error(ctx, err)
		return nil
	}

	return &User{
		root:        rsv,
		uid:         mutatedUser.UID,
		id:          store.ID(params.User),
		creation:    mutatedUser.Creation,
		displayName: mutatedUser.DisplayName,
		email:       mutatedUser.Email,
	}
}
