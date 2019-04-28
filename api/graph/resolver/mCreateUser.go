package resolver

import (
	"context"
	"demo/store/dbmod"

	"github.com/pkg/errors"
)

// CreateUser resolves Mutation.createUser
func (rsv *Resolver) CreateUser(
	ctx context.Context,
	params struct {
		Email       string
		DisplayName string
	},
) (*User, error) {
	newID, err := rsv.str.CreateUser(ctx, params.Email, params.DisplayName)
	if err != nil {
		return nil, err
	}

	var result struct {
		NewUser []dbmod.User `json:"newUser"`
	}
	if err := rsv.str.QueryVars(
		ctx,
		`query NewUser($id: string) {
			newUser(func: eq(User.id, $id)) {
				uid
				User.id
				User.creation
				User.email
				User.displayName
			}
		}`,
		map[string]string{
			"$id": string(newID),
		},
		&result,
	); err != nil {
		return nil, err
	}

	if len(result.NewUser) != 1 {
		return nil, errors.Errorf(
			"unexpected number of new users: %d",
			len(result.NewUser),
		)
	}

	newUser := result.NewUser[0]

	return &User{
		root:        rsv,
		uid:         *newUser.UID,
		id:          *newUser.ID,
		creation:    *newUser.Creation,
		displayName: *newUser.DisplayName,
		email:       *newUser.Email,
	}, nil
}
