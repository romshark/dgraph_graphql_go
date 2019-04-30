package dbmod

import (
	"time"

	"github.com/romshark/dgraph_graphql_go/store"
)

// Post defines the Post type query object
type Post struct {
	UID       string     `json:"uid"`
	ID        store.ID   `json:"Post.id"`
	Creation  time.Time  `json:"Post.creation"`
	Author    User       `json:"Post.author"`
	Title     string     `json:"Post.title"`
	Contents  string     `json:"Post.contents"`
	Reactions []Reaction `json:"Post.reaction"`
}
