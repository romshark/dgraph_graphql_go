package resolver

import (
	"context"
	"reflect"
	"time"

	"github.com/graph-gophers/graphql-go"
	"github.com/pkg/errors"
	"github.com/romshark/dgraph_graphql_go/store"
	"github.com/romshark/dgraph_graphql_go/store/dgraph"
	"github.com/romshark/dgraph_graphql_go/store/enum/emotion"
)

// Reaction represents the resolver of the identically named type
type Reaction struct {
	root       *Resolver
	uid        string
	authorUID  string
	subjectUID string
	id         store.ID
	creation   time.Time
	emotion    emotion.Emotion
	message    string
}

// ID resolves Reaction.id
func (rsv *Reaction) ID() store.ID {
	return rsv.id
}

// Creation resolves Reaction.creation
func (rsv *Reaction) Creation() graphql.Time {
	return graphql.Time{
		Time: rsv.creation,
	}
}

// Subject resolves Reaction.subject
func (rsv *Reaction) Subject(ctx context.Context) (*ReactionSubject, error) {
	var query struct {
		Subject []dgraph.ReactionSubject `json:"subject"`
	}
	if err := rsv.root.str.QueryVars(
		ctx,
		`query ReactionSubject($nodeId: string) {
			subject(func: uid($nodeId)) {
				uid

				Post.id
				Post.creation
				Post.author {
					uid
				}
				Post.title
				Post.contents

				Reaction.id
				Reaction.creation
				Reaction.author {
					uid
				}
				Reaction.message
				Reaction.emotion
			}
		}`,
		map[string]string{
			"$nodeId": rsv.subjectUID,
		},
		&query,
	); err != nil {
		rsv.root.error(ctx, err)
		return nil, err
	}

	subject := query.Subject[0]

	switch v := subject.V.(type) {
	case *dgraph.Post:
		return &ReactionSubject{&Post{
			root:      rsv.root,
			uid:       v.UID,
			id:        v.ID,
			creation:  v.Creation,
			title:     v.Title,
			contents:  v.Contents,
			authorUID: v.Author[0].UID,
		}}, nil
	case *dgraph.Reaction:
		return &ReactionSubject{&Reaction{
			root:       rsv.root,
			uid:        v.UID,
			authorUID:  v.Author[0].UID,
			subjectUID: rsv.uid,
			id:         v.ID,
			creation:   v.Creation,
			emotion:    v.Emotion,
			message:    v.Message,
		}}, nil
	}
	err := errors.Errorf(
		"unsupported union ReactionSubject type: %s",
		reflect.TypeOf(subject.V),
	)
	rsv.root.error(ctx, err)
	return nil, err
}

// Author resolves Reaction.author
func (rsv *Reaction) Author(ctx context.Context) (*User, error) {
	var query struct {
		Author []dgraph.User `json:"author"`
	}
	if err := rsv.root.str.QueryVars(
		ctx,
		`query ReactionAuthor($nodeId: string) {
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
		uid:         author.UID,
		id:          store.ID(author.ID),
		creation:    author.Creation,
		email:       author.Email,
		displayName: author.DisplayName,
	}, nil
}

// Emotion resolves Reaction.emotion
func (rsv *Reaction) Emotion() string {
	return string(rsv.emotion)
}

// Message resolves Reaction.message
func (rsv *Reaction) Message() string {
	return rsv.message
}

// Reactions resolves Reaction.reactions
func (rsv *Reaction) Reactions(ctx context.Context) ([]*Reaction, error) {
	var query struct {
		Reaction []dgraph.Reaction `json:"reaction"`
	}
	if err := rsv.root.str.QueryVars(
		ctx,
		`query SubReactions($nodeId: string) {
			reaction(func: uid($nodeId)) {
				Reaction.reactions {
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

	if len(query.Reaction) < 1 {
		return nil, nil
	}

	resolvers := make([]*Reaction, len(query.Reaction[0].Reactions))
	for i, subReaction := range query.Reaction[0].Reactions {
		resolvers[i] = &Reaction{
			root:       rsv.root,
			uid:        subReaction.UID,
			id:         subReaction.ID,
			authorUID:  subReaction.Author[0].UID,
			subjectUID: rsv.uid,
			creation:   subReaction.Creation,
			emotion:    subReaction.Emotion,
			message:    subReaction.Message,
		}
	}

	return resolvers, nil
}
