package store

import (
	"context"
	"encoding/json"
	"time"

	"github.com/dgraph-io/dgo/protos/api"
)

func newPost(
	ctx context.Context,
	txn transaction,
	ID ID,
	//TODO: authorUID must be a UID, not a string
	authorUID string,
	title string,
	contents string,
	creation time.Time,
) (err error) {
	// Marshal JSON
	var newPostJSON []byte
	newPostJSON, err = json.Marshal(struct {
		ID       string    `json:"Post.id"`
		Author   string    `json:"Post.author"`
		Title    string    `json:"Post.title"`
		Contents string    `json:"Post.contents"`
		Creation time.Time `json:"Post.creation"`
	}{
		Author:   authorUID,
		ID:       string(ID),
		Title:    title,
		Contents: contents,
		Creation: creation,
	})
	if err != nil {
		return
	}

	// Write
	_, err = txn.Mutation(ctx, &api.Mutation{
		SetJson: newPostJSON,
	})
	return
}
