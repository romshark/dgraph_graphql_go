package store

import (
	"context"
	"time"

	"github.com/pkg/errors"
	strerr "github.com/romshark/dgraph_graphql_go/store/errors"
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
	var res struct {
		ByID []struct {
			UID string `json:"uid"`
		} `json:"byId"`
		ByEmail []struct {
			UID string `json:"uid"`
		} `json:"byEmail"`
		ByDisplayName []struct {
			UID string `json:"uid"`
		} `json:"byDisplayName"`
	}
	err = txn.QueryVars(
		ctx,
		`query User(
			$id: string,
			$email: string,
			$displayName: string
		) {
			byId(func: eq(User.id, $id)) { uid }
			byEmail(func: eq(User.email, $email)) { uid }
			byDisplayName(func: eq(User.displayName, $displayName)) { uid }
		}`,
		map[string]string{
			"$email":       email,
			"$displayName": displayName,
		},
		&res,
	)
	if err != nil {
		return
	}

	if len(res.ByID) > 0 {
		err = errors.Errorf("duplicate User.id: %s", newID)
		return
	}
	if len(res.ByEmail) > 0 {
		err = strerr.Newf(
			strerr.ErrInvalidInput,
			"%d users with a similar email already exist",
			len(res.ByEmail),
		)
		return
	}
	if len(res.ByDisplayName) > 0 {
		err = strerr.Newf(
			strerr.ErrInvalidInput,
			"%d users with a similar displayName already exist",
			len(res.ByDisplayName),
		)
		return
	}

	err = newUser(ctx, txn, newID, email, displayName, time.Now())
	return
}
