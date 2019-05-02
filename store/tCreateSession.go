package store

import (
	"context"
	"encoding/json"
	"time"

	"github.com/dgraph-io/dgo/protos/api"
	strerr "github.com/romshark/dgraph_graphql_go/store/errors"
)

// CreateSession creates a new session and updates the indexes
func (str *store) CreateSession(
	ctx context.Context,
	email string,
	password string,
) (
	newUID UID,
	newKey string,
	creation time.Time,
	userId UID,
	err error,
) {
	// Validate inputs
	if len(email) < 1 || len(password) < 1 {
		err = strerr.New(strerr.ErrInvalidInput, "missing credentials")
		return
	}

	// Prepare
	newKey = str.sessionKeyGenerator.Generate()

	// Begin transaction
	txn, close := str.txn(&err)
	if err != nil {
		return
	}
	defer close()

	// Ensure user exists
	var res struct {
		ByEmail []struct {
			UID      string `json:"uid"`
			Password string `json:"User.password"`
		} `json:"byEmail"`
	}
	err = txn.QueryVars(
		ctx,
		`query User(
			$email: string
		) {
			byEmail(func: eq(User.email, $email)) {
				uid
				User.password
			}
		}`,
		map[string]string{
			"$email": email,
		},
		&res,
	)
	if err != nil {
		return
	}

	// Ensure the user exists and the password is correct
	if len(res.ByEmail) < 1 || !str.passwordHasher.Compare(
		[]byte(password),
		[]byte(res.ByEmail[0].Password),
	) {
		err = strerr.New(strerr.ErrWrongCreds, "wrong credentials")
		return
	}

	userId = UID{NodeID: res.ByEmail[0].UID}
	creation = time.Now()

	// Create new session
	var newSessionJSON []byte
	newSessionJSON, err = json.Marshal(struct {
		Key      string    `json:"Session.key"`
		Creation time.Time `json:"Session.creation"`
		User     UID       `json:"Session.user"`
	}{
		Key:      newKey,
		Creation: creation,
		User:     userId,
	})
	if err != nil {
		return
	}

	var sessCreationMut map[string]string
	sessCreationMut, err = txn.Mutation(ctx, &api.Mutation{
		SetJson: newSessionJSON,
	})
	if err != nil {
		return
	}
	newUID = UID{sessCreationMut["blank-0"]}

	// Update owner (User.sessions -> new session)
	updateOwner := struct {
		UID      string `json:"uid"`
		Sessions UID    `json:"User.sessions"`
	}{
		UID:      userId.NodeID,
		Sessions: newUID,
	}
	var updateOwnerJSON []byte
	updateOwnerJSON, err = json.Marshal(updateOwner)
	if err != nil {
		return
	}

	_, err = txn.Mutation(ctx, &api.Mutation{
		SetJson: updateOwnerJSON,
	})
	if err != nil {
		return
	}

	// Add the new session to the global Index
	var newSessionIndexJSON []byte
	newSessionIndexJSON, err = json.Marshal(struct {
		UID UID `json:"sessions"`
	}{
		UID: newUID,
	})
	if err != nil {
		return
	}

	_, err = txn.Mutation(ctx, &api.Mutation{
		SetJson: newSessionIndexJSON,
		Set:     nil,
	})

	return
}
