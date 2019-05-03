package dgraph

import (
	"time"
)

// User defines the User type query object
type User struct {
	UID                string     `json:"uid"`
	ID                 string     `json:"User.id"`
	Creation           time.Time  `json:"User.creation"`
	Email              string     `json:"User.email"`
	DisplayName        string     `json:"User.displayName"`
	Password           string     `json:"User.password"`
	Posts              []Post     `json:"User.posts"`
	Sessions           []Session  `json:"User.sessions"`
	PublishedReactions []Reaction `json:"User.publishedReactions"`
}
