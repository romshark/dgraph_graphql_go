package resolver

import (
	"context"
	"time"

	"github.com/graph-gophers/graphql-go"
	"github.com/romshark/dgraph_graphql_go/store"
	"github.com/romshark/dgraph_graphql_go/store/dgraph"
)

// Post represents the resolver of the identically named type
type Post struct {
	root      *Resolver
	uid       string
	authorUID string
	id        store.ID
	creation  time.Time
	title     string
	contents  string
}

// ID resolves Post.id
func (rsv *Post) ID() store.ID {
	return rsv.id
}

// Author resolves Post.author
func (rsv *Post) Author(ctx context.Context) (*User, error) {
	var query struct {
		Author []dgraph.User `json:"author"`
	}
	if err := rsv.root.str.QueryVars(
		ctx,
		`query Author($nodeId: string) {
			author(func: uid($nodeId)) {
				uid
				User.id
				User.creation
				User.email
				User.displayName
			}
		}`,
		map[string]string{
			"$nodeId": rsv.authorUID,
		},
		&query,
	); err != nil {
		rsv.root.error(ctx, err)
		return nil, err
	}

	author := query.Author[0]
	return &User{
		root:        rsv.root,
		uid:         rsv.authorUID,
		id:          store.ID(author.ID),
		creation:    author.Creation,
		email:       author.Email,
		displayName: author.DisplayName,
	}, nil
}

// Creation resolves Post.creation
func (rsv *Post) Creation() graphql.Time {
	return graphql.Time{
		Time: rsv.creation,
	}
}

// Title resolves Post.title
func (rsv *Post) Title() string {
	return rsv.title
}

// Contents resolves Post.contents
func (rsv *Post) Contents() string {
	return rsv.contents
}

// Reactions resolves Post.reactions
func (rsv *Post) Reactions(ctx context.Context) ([]*Reaction, error) {
	var query struct {
		Posts []dgraph.Post `json:"posts"`
	}
	if err := rsv.root.str.QueryVars(
		ctx,
		`query ReactionsToPost($nodeId: string) {
			posts(func: uid($nodeId)) {
				Post.reactions {
					uid
					Reaction.id
					Reaction.creation
					Reaction.emotion
					Reaction.message
					Reaction.author {
						uid
					}
				}
			}
		}`,
		map[string]string{
			"$nodeId": rsv.uid,
		},
		&query,
	); err != nil {
		rsv.root.error(ctx, err)
		return nil, err
	}

	if len(query.Posts) < 1 {
		return nil, nil
	}

	post := query.Posts[0]
	resolvers := make([]*Reaction, len(post.Reactions))
	for i, reaction := range post.Reactions {
		resolvers[i] = &Reaction{
			root:       rsv.root,
			uid:        reaction.UID,
			id:         reaction.ID,
			subjectUID: rsv.uid,
			authorUID:  reaction.Author[0].UID,
			creation:   reaction.Creation,
			emotion:    reaction.Emotion,
			message:    reaction.Message,
		}
	}

	return resolvers, nil
}
