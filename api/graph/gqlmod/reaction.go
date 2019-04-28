package gqlmod

import (
	"time"

	"github.com/romshark/dgraph_graphql_go/store"
	"github.com/romshark/dgraph_graphql_go/store/enum/emotion"
)

// Reaction defines the Reaction type query object
type Reaction struct {
	ID        *store.ID        `json:"id"`
	Subject   interface{}      `json:"subject"`
	Creation  *time.Time       `json:"creation"`
	Author    *User            `json:"author"`
	Message   *string          `json:"message"`
	Emotion   *emotion.Emotion `json:"emotion"`
	Reactions []Reaction       `json:"reactions"`
}
