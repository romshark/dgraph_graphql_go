package store

import (
	"context"
	"demo/store/enum/emotion"
	"encoding/json"
	"time"

	"github.com/dgraph-io/dgo/protos/api"
)

func newReaction(
	ctx context.Context,
	txn transaction,
	ID ID,
	//TODO: subjectUID must be a UID, not a string
	subjectUID string,
	//TODO: authorUID must be a UID, not a string
	authorUID string,
	emotion emotion.Emotion,
	message string,
	creation time.Time,
) (err error) {
	// Marshal JSON
	var newReactionJSON []byte
	newReactionJSON, err = json.Marshal(struct {
		ID       string    `json:"Post.id"`
		Subject  string    `json:"Post.subject"`
		Author   string    `json:"Post.author"`
		Emotion  string    `json:"Post.emotion"`
		Message  string    `json:"Post.message"`
		Creation time.Time `json:"Post.creation"`
	}{
		Subject:  subjectUID,
		Author:   authorUID,
		ID:       string(ID),
		Emotion:  string(emotion),
		Message:  message,
		Creation: creation,
	})
	if err != nil {
		return
	}

	// Write
	_, err = txn.Mutation(ctx, &api.Mutation{
		SetJson: newReactionJSON,
	})
	return
}
