package store

import (
	"context"
	"time"

	"github.com/pkg/errors"
)

// CreateUser creates a new user account
func (str *store) CreateUser(
	ctx context.Context,
	email string,
	displayName string,
) (newID ID, err error) {
	// Validate inputs
	if err := ValidateUserDisplayName(displayName); err != nil {
		return "", err
	}

	// Prepare
	newID = NewID()

	// Begin transaction
	txn, close := str.txn(&err)
	if err != nil {
		return
	}
	defer close()

	// Ensure no users with a similar email already exist
	//TODO: check ID and displayName as well
	var res struct {
		Users []string `json:"users"`
	}
	err = txn.QueryVars(
		ctx,
		`query User($email: string) {
			users(func: eq(User.email, $email)) { uid }
		}`,
		map[string]string{
			"$email": email,
		},
		&res,
	)
	if err != nil {
		return
	}

	if len(res.Users) > 0 {
		err = errors.Errorf(
			"%d users with a similar email already exist",
			len(res.Users),
		)
		return
	}

	err = newUser(ctx, txn, newID, email, displayName, time.Now())
	return
}
