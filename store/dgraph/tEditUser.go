package dgraph

import (
	"context"
	"encoding/json"

	"github.com/dgraph-io/dgo/v2/protos/api"
	"github.com/romshark/dgraph_graphql_go/store"
	strerr "github.com/romshark/dgraph_graphql_go/store/errors"
)

// EditUser edits an existing user profile
func (str *impl) EditUser(
	ctx context.Context,
	user store.ID,
	editor store.ID,
	newEmail *string,
	newPassword *string,
) (
	result store.User,
	changes struct {
		Email    bool
		Password bool
	},
	err error,
) {
	result.ID = user

	// Begin transaction
	txn, close := str.txn(&err)
	if err != nil {
		return
	}
	defer close()

	// Ensure user and editor exist
	var qr struct {
		User   []User `json:"user"`
		Editor []User `json:"editor"`
	}
	err = txn.QueryVars(
		ctx,
		`query User(
			$id: string,
			$editorId: string
		) {
			user(func: eq(User.id, $id)) {
				uid
				User.creation
				User.displayName
				User.email
				User.password
			}
			editor(func: eq(User.id, $editorId)) { uid }
		}`,
		map[string]string{
			"$id":       string(user),
			"$editorId": string(editor),
		},
		&qr,
	)
	if err != nil {
		return
	}

	if len(qr.User) < 1 {
		err = strerr.New(strerr.ErrInvalidInput, "user profile not found")
		return
	}
	if len(qr.Editor) < 1 {
		err = strerr.Newf(strerr.ErrInvalidInput, "editor not found")
		return
	}

	if newEmail != nil {
		result.Email = *newEmail
		if qr.User[0].Email == *newEmail {
			newEmail = nil
		} else {
			changes.Email = true
		}
	} else {
		result.Email = qr.User[0].Email
	}
	if newPassword != nil {
		result.Password = *newPassword
		if qr.User[0].Password == *newPassword {
			newPassword = nil
		} else {
			changes.Password = true
		}
	} else {
		result.Password = qr.User[0].Password
	}

	result.UID = qr.User[0].UID
	result.Creation = qr.User[0].Creation
	result.DisplayName = qr.User[0].DisplayName

	// Edit the user profile
	var mutatedUserJSON []byte
	mutatedUserJSON, err = json.Marshal(struct {
		UID         string  `json:"uid"`
		NewEmail    *string `json:"User.email,omitempty"`
		NewPassword *string `json:"User.password,omitempty"`
	}{
		UID:         result.UID,
		NewEmail:    newEmail,
		NewPassword: newPassword,
	})
	if err != nil {
		return
	}
	_, err = txn.Mutation(ctx, &api.Mutation{
		SetJson: mutatedUserJSON,
	})
	if err != nil {
		return
	}

	return
}
