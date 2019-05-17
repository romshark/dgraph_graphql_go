package store

import (
	"time"
)

// User represents a User entity
type User struct {
	GraphNode

	ID                 ID
	Creation           time.Time
	Email              string
	DisplayName        string
	Password           string
	Posts              []Post
	Sessions           []Session
	PublishedReactions []Reaction
}
