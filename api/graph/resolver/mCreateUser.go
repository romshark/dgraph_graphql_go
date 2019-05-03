package resolver

import (
	"context"
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
		uid:         transactRes.UID.NodeID,
		id:          transactRes.ID,
		creation:    transactRes.CreationTime,
		displayName: params.DisplayName,
		email:       params.Email,
	}, nil
}
