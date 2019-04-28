package store

import (
	"context"

	"github.com/dgraph-io/dgo/protos/api"
)

func (str *store) setupSchema(ctx context.Context) error {
	// TODO: avoid DropAll
	if err := str.db.Alter(ctx, &api.Operation{
		DropAll: true,
	}); err != nil {
		panic(err)
	}

	return str.db.Alter(ctx, &api.Operation{
		Schema: `
			User.id: string @index(exact) .
			User.creation: dateTime .
			User.email: string @index(exact) .
			User.displayName: string @index(exact) .
			User.posts: uid .

			Post.id: string @index(exact) .
			Post.creation: dateTime .
			Post.title: string .
			Post.contents: string .
			Post.author: uid .
			Post.reactions: uid .

			Reaction.id: string @index(exact) .
			Reaction.creation: dateTime .
			Reaction.emotion: string .
			Reaction.message: string .
			Reaction.author: uid .
			Reaction.reactions: uid .
		`,
	})
}
