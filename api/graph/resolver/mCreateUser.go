package resolver

import (
	"context"

	"github.com/pkg/errors"
	"github.com/romshark/dgraph_graphql_go/store/dbmod"
)

// CreateUser resolves Mutation.createUser
func (rsv *Resolver) CreateUser(
	ctx context.Context,
	params struct {
		Email       string
		DisplayName string
	},
) (*User, error) {
	newUID, newID, err := rsv.str.CreateUser(
		ctx,
		params.Email,
		params.DisplayName,
	)
	if err != nil {
		rsv.error(ctx, err)
		return nil, err
	}

	var result struct {
		NewUser []dbmod.User `json:"newUser"`
	}
	if err := rsv.str.QueryVars(
		ctx,
		`query NewUser($nodeId: string) {
			newUser(func: uid($nodeId)) {
				User.creation
				User.email
				User.displayName
			}
		}`,
		map[string]string{
			"$nodeId": newUID.NodeID,
		},
		&result,
	); err != nil {
		rsv.error(ctx, err)
		return nil, err
	}

	if len(result.NewUser) != 1 {
		err := errors.Errorf(
			"unexpected number of new users: %d",
			len(result.NewUser),
		)
		rsv.error(ctx, err)
		return nil, err
	}

	newUser := result.NewUser[0]

	return &User{
		root:        rsv,
		uid:         newUID,
		id:          newID,
		creation:    newUser.Creation,
		displayName: newUser.DisplayName,
		email:       newUser.Email,
	}, nil
}
