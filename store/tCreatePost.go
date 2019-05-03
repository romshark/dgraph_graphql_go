package store

import (
	"context"
	"encoding/json"
	"time"

	"github.com/dgraph-io/dgo/protos/api"
	"github.com/pkg/errors"
	strerr "github.com/romshark/dgraph_graphql_go/store/errors"
)

// CreatePost creates a new post
func (str *store) CreatePost(
	ctx context.Context,
	authorID ID,
	title string,
	contents string,
) (
	result struct {
		UID          UID
		ID           ID
		AuthorUID    UID
		CreationTime time.Time
	},
	err error,
) {
	// Validate input
	err = ValidatePostTitle(title)
	if err != nil {
		err = strerr.Wrap(strerr.ErrInvalidInput, err)
		return
	}
	err = ValidatePostContents(contents)
	if err != nil {
		err = strerr.Wrap(strerr.ErrInvalidInput, err)
		return
	}

	// Prepare
	result.ID = NewID()
	result.CreationTime = time.Now()

	// Begin transaction
	txn, close := str.txn(&err)
	if err != nil {
		return
	}
	defer close()

	// Ensure author exists
	var res struct {
		ByID   []UID `json:"byId"`
		Author []UID `json:"author"`
	}
	err = txn.QueryVars(
		ctx,
		`query User(
			$id: string,
			$authorId: string
		) {
			byId(func: eq(Post.id, $id)) { uid }
			author(func: eq(User.id, $authorId)) { uid }
		}`,
		map[string]string{
			"$id":       string(result.ID),
			"$authorId": string(authorID),
		},
		&res,
	)
	if err != nil {
		return
	}

	if len(res.ByID) > 0 {
		err = errors.Errorf("duplicate Post.id: %s", result.ID)
		return
	}
	if len(res.Author) < 1 {
		err = strerr.Newf(
			strerr.ErrInvalidInput,
			"author not found",
		)
		return
	}

	result.AuthorUID = res.Author[0]

	// Create new post
	var newPostJSON []byte
	newPostJSON, err = json.Marshal(struct {
		ID       string    `json:"Post.id"`
		Author   UID       `json:"Post.author"`
		Title    string    `json:"Post.title"`
		Contents string    `json:"Post.contents"`
		Creation time.Time `json:"Post.creation"`
	}{
		Author:   result.AuthorUID,
		ID:       string(result.ID),
		Title:    title,
		Contents: contents,
		Creation: result.CreationTime,
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
	result.UID = UID{postCreationMut["blank-0"]}

	// Update author (User.posts -> new post)
	var updatedAuthorJSON []byte
	updatedAuthorJSON, err = json.Marshal(struct {
		UID   string `json:"uid"`
		Posts UID    `json:"User.posts"`
	}{
		UID:   result.AuthorUID.NodeID,
		Posts: result.UID,
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

	// Add the new post to the global Index
	var newPostsIndexJSON []byte
	newPostsIndexJSON, err = json.Marshal(struct {
		UID UID `json:"posts"`
	}{
		UID: result.UID,
	})
	if err != nil {
		return
	}

	_, err = txn.Mutation(ctx, &api.Mutation{
		SetJson: newPostsIndexJSON,
	})

	return
}
