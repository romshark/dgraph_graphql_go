package dgraph

import (
	"context"
	"encoding/json"

	"github.com/dgraph-io/dgo/protos/api"
	"github.com/pkg/errors"
	"github.com/romshark/dgraph_graphql_go/api/graph/auth"
	"github.com/romshark/dgraph_graphql_go/store"
	strerr "github.com/romshark/dgraph_graphql_go/store/errors"
)

// EditPost edits an existing post
func (str *impl) EditPost(
	ctx context.Context,
	post store.ID,
	editor store.ID,
	newTitle *string,
	newContents *string,
) (
	result store.Post,
	changes struct {
		Title    bool
		Contents bool
	},
	err error,
) {
	result.ID = post

	// Begin transaction
	txn, close := str.txn(&err)
	if err != nil {
		return
	}
	defer close()

	// Ensure post and editor exist
	var qr struct {
		Post   []Post `json:"post"`
		Editor []User `json:"editor"`
	}
	err = txn.QueryVars(
		ctx,
		`query User(
			$id: string,
			$editorId: string
		) {
			post(func: eq(Post.id, $id)) {
				uid
				Post.author {
					uid
					User.id
				}
				Post.creation
				Post.title
				Post.contents
			}
			editor(func: eq(User.id, $editorId)) { uid }
		}`,
		map[string]string{
			"$id":       string(post),
			"$editorId": string(editor),
		},
		&qr,
	)
	if err != nil {
		return
	}

	if len(qr.Post) < 1 {
		err = errors.New("post not found")
		return
	}
	if len(qr.Editor) < 1 {
		err = strerr.Newf(strerr.ErrInvalidInput, "editor not found")
		return
	}

	// Check permission
	if err = auth.Authorize(ctx, auth.IsOwner{
		Owner: store.ID(qr.Post[0].Author[0].ID),
	}); err != nil {
		return
	}

	if newTitle != nil {
		result.Title = *newTitle
		if qr.Post[0].Title == *newTitle {
			newTitle = nil
		} else {
			changes.Title = true
		}
	} else {
		result.Title = qr.Post[0].Title
	}
	if newContents != nil {
		result.Contents = *newContents
		if qr.Post[0].Contents == *newContents {
			newContents = nil
		} else {
			changes.Contents = true
		}
	} else {
		result.Contents = qr.Post[0].Contents
	}

	result.UID = qr.Post[0].UID
	result.Creation = qr.Post[0].Creation
	result.Author = &store.User{
		GraphNode: store.GraphNode{
			UID: qr.Post[0].Author[0].UID,
		},
	}

	// Edit the post
	var mutatedPostJSON []byte
	mutatedPostJSON, err = json.Marshal(struct {
		UID         string  `json:"uid"`
		NewTitle    *string `json:"Post.title,omitempty"`
		NewContents *string `json:"Post.contents,omitempty"`
	}{
		UID:         result.UID,
		NewTitle:    newTitle,
		NewContents: newContents,
	})
	if err != nil {
		return
	}
	_, err = txn.Mutation(ctx, &api.Mutation{
		SetJson: mutatedPostJSON,
	})
	if err != nil {
		return
	}

	return
}
