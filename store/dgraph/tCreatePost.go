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

// CreatePost creates a new post
func (str *impl) CreatePost(
	ctx context.Context,
	creationTime time.Time,
	authorID store.ID,
	title string,
	contents string,
) (
	result store.Post,
	err error,
) {
	result.Title = title
	result.Contents = contents
	result.Creation = creationTime

	// Prepare
	result.ID = store.NewID()

	// Begin transaction
	txn, close := str.txn(&err)
	if err != nil {
		return
	}
	defer close()

	// Get author and posts list meta information
	var qr struct {
		ByID         []UID                   `json:"byId"`
		Author       []UID                   `json:"author"`
		PostsCount   []struct{ Count int32 } `json:"postsCount"`
		PostsVersion []UID                   `json:"postsVersion"`
	}
	err = txn.QueryVars(
		ctx,
		`query User(
			$id: string,
			$authorId: string
		) {
			byId(func: eq(Post.id, $id)) { uid }
			author(func: eq(User.id, $authorId)) { uid }
			postsCount(func: has(<posts>)) {
				count: count(uid)
			}
			postsVersion(func: has(posts.version)) {
				uid
			}
		}`,
		map[string]string{
			"$id":       string(result.ID),
			"$authorId": string(authorID),
		},
		&qr,
	)
	if err != nil {
		return
	}

	if len(qr.ByID) > 0 {
		err = errors.Errorf("duplicate Post.id: %s", result.ID)
		return
	}
	if len(qr.Author) < 1 {
		err = strerr.Newf(
			strerr.ErrInvalidInput,
			"author not found",
		)
		return
	}

	result.Author = &store.User{
		GraphNode: store.GraphNode{
			UID: qr.Author[0].NodeID,
		},
	}

	// Create new post
	var newPostJSON []byte
	newPostJSON, err = json.Marshal(struct {
		ID       string    `json:"Post.id"`
		Author   UID       `json:"Post.author"`
		Title    string    `json:"Post.title"`
		Contents string    `json:"Post.contents"`
		Creation time.Time `json:"Post.creation"`
	}{
		Author:   UID{NodeID: result.Author.UID},
		ID:       string(result.ID),
		Title:    title,
		Contents: contents,
		Creation: creationTime,
	})
	if err != nil {
		return
	}
	var postCreationMut map[string]string
	postCreationMut, err = txn.Mutation(ctx, &api.Mutation{
		SetJson: newPostJSON,
	})
	if err != nil {
		return
	}
	result.UID = postCreationMut["blank-0"]

	// Update author (User.posts -> new post)
	var updatedAuthorJSON []byte
	updatedAuthorJSON, err = json.Marshal(struct {
		UID   string `json:"uid"`
		Posts UID    `json:"User.posts"`
	}{
		UID:   result.Author.UID,
		Posts: UID{NodeID: result.UID},
	})
	if err != nil {
		return
	}

	_, err = txn.Mutation(ctx, &api.Mutation{
		SetJson: updatedAuthorJSON,
	})
	if err != nil {
		return
	}

	// Add the new post to the global index
	var updateJSON []byte
	updateJSON, err = json.Marshal(struct {
		UID UID `json:"posts"`
	}{
		UID: UID{NodeID: result.UID},
	})
	if err != nil {
		return
	}

	_, err = txn.Mutation(ctx, &api.Mutation{
		SetJson: updateJSON,
	})
	if err != nil {
		return
	}

	// Update posts.version
	updateJSON, err = json.Marshal(struct {
		UID     string `json:"uid"`
		Version string `json:"posts.version"`
	}{
		UID:     qr.PostsVersion[0].NodeID,
		Version: string(store.NewID()),
	})
	if err != nil {
		return
	}

	_, err = txn.Mutation(ctx, &api.Mutation{
		SetJson: updateJSON,
	})
	if err != nil {
		return
	}

	return
}
