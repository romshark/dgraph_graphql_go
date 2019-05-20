package store

import (
	"context"
	"time"

	"github.com/romshark/dgraph_graphql_go/store/enum/emotion"
)

// MutableStore interfaces a transactional store
type MutableStore interface {
	CreateSession(
		ctx context.Context,
		key string,
		creationTime time.Time,
		email string,
		password string,
	) (
		result Session,
		err error,
	)

	CloseSession(
		ctx context.Context,
		key string,
	) (
		result bool,
		err error,
	)

	CloseAllSessions(
		ctx context.Context,
		user ID,
	) (
		result []string,
		err error,
	)

	CreatePost(
		ctx context.Context,
		creationTime time.Time,
		author ID,
		title string,
		contents string,
	) (
		result Post,
		err error,
	)

	CreateReaction(
		ctx context.Context,
		creationTime time.Time,
		author ID,
		subject ID,
		emotion emotion.Emotion,
		message string,
	) (
		result Reaction,
		err error,
	)

	CreateUser(
		ctx context.Context,
		creationTime time.Time,
		email string,
		displayName string,
		passwordHash string,
	) (
		result User,
		err error,
	)

	EditPost(
		ctx context.Context,
		post ID,
		editor ID,
		newTitle *string,
		newContents *string,
	) (
		result Post,
		changes struct {
			Title    bool
			Contents bool
		},
		err error,
	)

	EditUser(
		ctx context.Context,
		user ID,
		editor ID,
		newEmail *string,
		newPassword *string,
	) (
		result User,
		changes struct {
			Email    bool
			Password bool
		},
		err error,
	)

	EditReaction(
		ctx context.Context,
		reaction ID,
		editor ID,
		newMessage string,
	) (
		result Reaction,
		changes struct {
			Message bool
		},
		err error,
	)
}

// Store interfaces a store implementation
type Store interface {
	Prepare() error

	MutableStore

	Query(
		ctx context.Context,
		query string,
		result interface{},
	) error

	QueryVars(
		ctx context.Context,
		query string,
		vars map[string]string,
		result interface{},
	) error
}
