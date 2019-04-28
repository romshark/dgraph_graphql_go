package dbmod

import (
	"time"

	"github.com/romshark/dgraph_graphql_go/store"
	"github.com/romshark/dgraph_graphql_go/store/enum/emotion"
)

// Reaction defines the Reaction type query object
type Reaction struct {
	UID       *string          `json:"uid"`
	ID        *store.ID        `json:"Reaction.id"`
	Subject   interface{}      `json:"Reaction.subject"`
	Creation  *time.Time       `json:"Reaction.creation"`
	Author    *User            `json:"Reaction.author"`
	Message   *string          `json:"Reaction.message"`
	Emotion   *emotion.Emotion `json:"Reaction.emotion"`
	Reactions []Reaction       `json:"Reaction.reactions"`
}
