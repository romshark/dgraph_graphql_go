package store

import (
	"context"
	"encoding/json"
	"time"

	"github.com/dgraph-io/dgo/protos/api"
	"github.com/pkg/errors"
	strerr "github.com/romshark/dgraph_graphql_go/store/errors"
)

// CreateUser creates a new user account and adds it to the global index
func (str *store) CreateUser(
	ctx context.Context,
	email string,
	displayName string,
) (newUID UID, newID ID, err error) {
	// Validate inputs
	if err := ValidateUserDisplayName(displayName); err != nil {
		return UID{}, "", strerr.Wrap(strerr.ErrInvalidInput, err)
	}
	if err := ValidateEmail(email); err != nil {
		return UID{}, "", strerr.Wrap(strerr.ErrInvalidInput, err)
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

	// Create user account
	var newUserJSON []byte
	newUserJSON, err = json.Marshal(struct {
		ID          string    `json:"User.id"`
		Email       string    `json:"User.email"`
		DisplayName string    `json:"User.displayName"`
		Creation    time.Time `json:"User.creation"`
	}{
		ID:          string(newID),
		Email:       email,
		DisplayName: displayName,
		Creation:    time.Now(),
	})
	if err != nil {
		return
	}

	var userCreationMut map[string]string
	userCreationMut, err = txn.Mutation(ctx, &api.Mutation{
		SetJson: newUserJSON,
	})
	if err != nil {
		return
	}
	newUID = UID{userCreationMut["blank-0"]}

	// Add the new account to the global Index
	var newUsersIndexJSON []byte
	newUsersIndexJSON, err = json.Marshal(struct {
		UID UID `json:"users"`
	}{
		UID: newUID,
	})
	if err != nil {
		return
	}

	_, err = txn.Mutation(ctx, &api.Mutation{
		SetJson: newUsersIndexJSON,
		Set:     nil,
	})

	return
}
