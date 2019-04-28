package resolver

import (
	"demo/store"
	"demo/store/enum/emotion"
	"time"

	"github.com/graph-gophers/graphql-go"
)

// Reaction represents the resolver of the identically named type
type Reaction struct {
	root     *Resolver
	uid      string
	id       store.ID
	creation time.Time
	emotion  emotion.Emotion
	message  string
}

// Id resolves Reaction.id
func (rsv *Reaction) Id() store.ID {
	return rsv.id
}

// Creation resolves Reaction.creation
func (rsv *Reaction) Creation() graphql.Time {
	return graphql.Time{
		Time: rsv.creation,
	}
}

// Subject resolves Reaction.subject
func (rsv *Reaction) Subject() *ReactionSubject {
	return nil
}

// Author resolves Reaction.author
func (rsv *Reaction) Author() *User {
	return nil
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
func (rsv *Reaction) Reactions() ([]*Reaction, error) {
	return nil, nil
}
