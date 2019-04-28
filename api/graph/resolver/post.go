package resolver

import (
	"context"
	"demo/store"
	"demo/store/dbmod"
	"time"

	"github.com/graph-gophers/graphql-go"
)

// Post represents the resolver of the identically named type
type Post struct {
	root     *Resolver
	uid      string
	id       store.ID
	creation time.Time
	title    string
	contents string
}

// Id resolves Post.id
func (rsv *Post) Id() store.ID {
	return rsv.id
}

// Author resolves Post.author
func (rsv *Post) Author(ctx context.Context) (*User, error) {
	var query struct {
		Author dbmod.User `json:"author"`
	}
	if err := rsv.root.str.QueryVars(
		ctx,
		`query Author($uid: uid) {
			author(func: uid($uid)) {
				Post.author {
					uid
					User.id
					User.creation
					User.email
					User.displayName
					User.posts
				}
			}
		}`,
		map[string]string{
			"$uid": rsv.uid,
		},
		&query,
	); err != nil {
		return nil, err
	}
	return &User{
		root:        rsv.root,
		uid:         *query.Author.UID,
		id:          *query.Author.ID,
		creation:    *query.Author.Creation,
		email:       *query.Author.Email,
		displayName: *query.Author.DisplayName,
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
func (rsv *Post) Reactions() ([]*Reaction, error) {
	return nil, nil
}
