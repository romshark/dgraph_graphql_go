package store

import (
	"context"
	"encoding/json"
	"time"

	"github.com/dgraph-io/dgo/protos/api"
)

func newUser(
	ctx context.Context,
	txn transaction,
	ID ID,
	email string,
	displayName string,
	creation time.Time,
) (err error) {
	// Marshal JSON
	var newUserJSON []byte
	newUserJSON, err = json.Marshal(struct {
		ID          string    `json:"User.id"`
		Email       string    `json:"User.email"`
		DisplayName string    `json:"User.displayName"`
		Creation    time.Time `json:"User.creation"`
	}{
		ID:          string(ID),
		Email:       email,
		DisplayName: displayName,
		Creation:    creation,
	})
	if err != nil {
		return
	}

	// Write
	_, err = txn.Mutation(ctx, &api.Mutation{
		SetJson: newUserJSON,
	})
	return
}
