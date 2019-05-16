package dgraph

import (
	"context"
	"encoding/json"
	"time"

	"github.com/dgraph-io/dgo/protos/api"
	"github.com/romshark/dgraph_graphql_go/store"
	strerr "github.com/romshark/dgraph_graphql_go/store/errors"
)

// CreateSession creates a new session and updates the indexes
func (str *impl) CreateSession(
	ctx context.Context,
	key string,
	creation time.Time,
	email string,
	password string,
) (
	result struct {
		UID     string
		UserID  store.ID
		UserUID string
	},
	err error,
) {
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
			ID       string `json:"User.id"`
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
				User.id
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
	if len(res.ByEmail) < 1 || !str.comparePassword(
		password,
		res.ByEmail[0].Password,
	) {
		err = strerr.New(strerr.ErrWrongCreds, "wrong credentials")
		return
	}

	result.UserID = store.ID(res.ByEmail[0].ID)
	result.UserUID = res.ByEmail[0].UID

	// Create new session
	var newSessionJSON []byte
	newSessionJSON, err = json.Marshal(struct {
		Key      string    `json:"Session.key"`
		Creation time.Time `json:"Session.creation"`
		User     UID       `json:"Session.user"`
	}{
		Key:      key,
		Creation: creation,
		User:     UID{NodeID: result.UserUID},
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
	result.UID = sessCreationMut["blank-0"]

	// Update owner (User.sessions -> new session)
	var updateOwnerJSON []byte
	updateOwnerJSON, err = json.Marshal(struct {
		UID      string `json:"uid"`
		Sessions UID    `json:"User.sessions"`
	}{
		UID:      result.UserUID,
		Sessions: UID{NodeID: result.UID},
	})
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
		UID: UID{NodeID: result.UID},
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
