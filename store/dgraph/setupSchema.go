package dgraph

import (
	"context"

	"github.com/dgraph-io/dgo/protos/api"
)

func (str *impl) setupSchema(ctx context.Context) error {
	return str.db.Alter(ctx, &api.Operation{
		Schema: `
			users: uid .
			posts: uid .
			sessions: uid @reverse .

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
	})
}
