package dbmod

import (
	"time"

	"github.com/romshark/dgraph_graphql_go/store"
)

// User defines the User type query object
type User struct {
	UID                string     `json:"uid"`
	ID                 store.ID   `json:"User.id"`
	Creation           time.Time  `json:"User.creation"`
	Email              string     `json:"User.email"`
	DisplayName        string     `json:"User.displayName"`
	Password           string     `json:"User.password"`
	Posts              []Post     `json:"User.posts"`
	Sessions           []Session  `json:"User.sessions"`
	PublishedReactions []Reaction `json:"User.publishedReactions"`
}
