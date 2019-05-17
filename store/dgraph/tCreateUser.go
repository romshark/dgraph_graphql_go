package dgraph

import (
	"context"
	"encoding/json"
	"time"

	"github.com/dgraph-io/dgo/protos/api"
	"github.com/pkg/errors"
	"github.com/romshark/dgraph_graphql_go/store"
	strerr "github.com/romshark/dgraph_graphql_go/store/errors"
)

// CreateUser creates a new user account and adds it to the global index
func (str *impl) CreateUser(
	ctx context.Context,
	creationTime time.Time,
	email string,
	displayName string,
	passwordHash string,
) (
	result store.User,
	err error,
) {
	result.Creation = creationTime
	result.Email = email
	result.DisplayName = displayName
	result.Password = passwordHash

	// Prepare
	result.ID = store.NewID()

	// Begin transaction
	txn, close := str.txn(&err)
	if err != nil {
		return
	}
	defer close()

	// Ensure no users with a similar email already exist
	var qr struct {
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
		&qr,
	)
	if err != nil {
		return
	}

	if len(qr.ByID) > 0 {
		err = errors.Errorf("duplicate User.id: %s", result.ID)
		return
	}
	if len(qr.ByEmail) > 0 {
		err = strerr.Newf(
			strerr.ErrInvalidInput,
			"%d users with a similar email already exist",
			len(qr.ByEmail),
		)
		return
	}
	if len(qr.ByDisplayName) > 0 {
		err = strerr.Newf(
			strerr.ErrInvalidInput,
			"%d users with a similar displayName already exist",
			len(qr.ByDisplayName),
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
		Password    string    `json:"User.password"`
	}{
		ID:          string(result.ID),
		Email:       email,
		DisplayName: displayName,
		Creation:    creationTime,
		Password:    string(passwordHash),
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
	result.UID = userCreationMut["blank-0"]

	// Add the new account to the global Index
	var newUsersIndexJSON []byte
	newUsersIndexJSON, err = json.Marshal(struct {
		UID UID `json:"users"`
	}{
		UID: UID{NodeID: result.UID},
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
