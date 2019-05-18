package dgraph

import (
	"context"
	"encoding/json"

	"github.com/dgraph-io/dgo/protos/api"
	"github.com/romshark/dgraph_graphql_go/api/graph/auth"
	"github.com/romshark/dgraph_graphql_go/store"
	strerr "github.com/romshark/dgraph_graphql_go/store/errors"
)

// EditReaction edits an existing reaction
func (str *impl) EditReaction(
	ctx context.Context,
	reaction store.ID,
	editor store.ID,
	newMessage string,
) (
	result store.Reaction,
	changes struct {
		Message bool
	},
	err error,
) {
	result.ID = reaction

	// Begin transaction
	txn, close := str.txn(&err)
	if err != nil {
		return
	}
	defer close()

	// Ensure reaction and editor exist
	var qr struct {
		Reaction []Reaction `json:"reaction"`
		Editor   []User     `json:"editor"`
	}
	err = txn.QueryVars(
		ctx,
		`query Reaction(
			$id: string,
			$editorId: string
		) {
			reaction(func: eq(Reaction.id, $id)) {
				uid
				Reaction.subject {
					uid
					Post.id
					Reaction.id
				}
				Reaction.author {
					uid
					User.id
				}
				Reaction.creation
				Reaction.emotion
				Reaction.message
			}
			editor(func: eq(User.id, $editorId)) { uid }
		}`,
		map[string]string{
			"$id":       string(reaction),
			"$editorId": string(editor),
		},
		&qr,
	)
	if err != nil {
		return
	}

	if len(qr.Reaction) < 1 {
		err = strerr.New(strerr.ErrInvalidInput, "reaction not found")
		return
	}
	if len(qr.Editor) < 1 {
		err = strerr.Newf(strerr.ErrInvalidInput, "editor not found")
		return
	}

	// Check permission
	if err = auth.Authorize(ctx, auth.IsOwner{
		Owner: store.ID(qr.Reaction[0].Author[0].ID),
	}); err != nil {
		return
	}

	result.Message = newMessage
	if qr.Reaction[0].Message != newMessage {
		changes.Message = true
	}

	react := qr.Reaction[0]

	result.UID = react.UID
	result.Creation = react.Creation
	result.Author = &store.User{
		GraphNode: store.GraphNode{
			UID: react.Author[0].UID,
		},
	}
	result.Emotion = react.Emotion

	switch subject := react.Subject[0].V.(type) {
	case *Post:
		result.Subject = store.Post{
			GraphNode: store.GraphNode{
				UID: subject.UID,
			},
		}
	case *Reaction:
		result.Subject = store.Reaction{
			GraphNode: store.GraphNode{
				UID: subject.UID,
			},
		}
	}

	// Edit the reaction
	var mutatedReactionJSON []byte
	mutatedReactionJSON, err = json.Marshal(struct {
		UID        string `json:"uid"`
		NewMessage string `json:"Reaction.message"`
	}{
		UID:        result.UID,
		NewMessage: newMessage,
	})
	if err != nil {
		return
	}
	_, err = txn.Mutation(ctx, &api.Mutation{
		SetJson: mutatedReactionJSON,
	})
	if err != nil {
		return
	}

	return
}
