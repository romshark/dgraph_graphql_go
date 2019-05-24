package dgraph

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/dgraph-io/dgo/protos/api"
	"github.com/romshark/dgraph_graphql_go/store"
)

func (str *impl) setupSchema(ctx context.Context) (err error) {
	if err = str.db.Alter(ctx, &api.Operation{
		Schema: `
			users: uid .
			posts: uid @count .
			sessions: uid @reverse .

			posts.version: string .

			Session.key: string @index(exact) .
			Session.creation: dateTime .
			Session.user: uid .

			User.id: string @index(exact) .
			User.creation: dateTime .
			User.email: string @index(exact) .
			User.displayName: string @index(exact) .
			User.posts: uid .
			User.sessions: uid .
			User.password: string .
			User.publishedReactions: uid .

			Post.id: string @index(exact) .
			Post.creation: dateTime .
			Post.title: string .
			Post.contents: string .
			Post.author: uid .
			Post.reactions: uid .

			Reaction.id: string @index(exact) .
			Reaction.subject: uid .
			Reaction.creation: dateTime .
			Reaction.emotion: string .
			Reaction.message: string .
			Reaction.author: uid .
			Reaction.reactions: uid .
		`,
	}); err != nil {
		return fmt.Errorf("schema setup: %s", err)
	}

	// Begin setup transaction
	txn, close := str.txn(&err)
	if err != nil {
		return
	}
	defer close()

	var qr struct {
		PostsVersion []UID `json:"postsVersion"`
	}
	if err = txn.Query(
		ctx,
		`{
			postsVersion(func: has(posts.version)) {
				uid
			}
		}`,
		&qr,
	); err != nil {
		return
	}

	if len(qr.PostsVersion) < 1 {
		// Define the first posts.version entry
		var addJSON []byte
		addJSON, err = json.Marshal(struct {
			PostsVersion string `json:"posts.version"`
		}{
			PostsVersion: string(store.NewID()),
		})
		if err != nil {
			return
		}

		_, err = txn.Mutation(ctx, &api.Mutation{
			SetJson: addJSON,
		})
		if err != nil {
			return
		}
	}

	return nil
}
