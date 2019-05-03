package gqlmod

import (
	"time"

	"github.com/romshark/dgraph_graphql_go/store"
)

// Post defines the Post type query object
type Post struct {
	ID        *store.ID  `json:"id"`
	Creation  *time.Time `json:"creation"`
	Author    *User      `json:"author"`
	Title     *string    `json:"title"`
	Contents  *string    `json:"contents"`
	Reactions []Reaction `json:"reactions"`
}
