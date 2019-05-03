package gqlmod

import (
	"time"

	"github.com/romshark/dgraph_graphql_go/store"
)

// User defines the User type query object
type User struct {
	ID                 *store.ID  `json:"id"`
	Creation           *time.Time `json:"creation"`
	Email              *string    `json:"email"`
	DisplayName        *string    `json:"displayName"`
	Posts              []Post     `json:"posts"`
	Sessions           []Session  `json:"sessions"`
	PublishedReactions []Reaction `json:"publishedReactions"`
}
