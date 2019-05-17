package dgraph

import (
	"time"

	"github.com/romshark/dgraph_graphql_go/store"
	"github.com/romshark/dgraph_graphql_go/store/enum/emotion"
)

// Reaction represents a database model for the Reaction entity
type Reaction struct {
	UID       string            `json:"uid"`
	ID        store.ID          `json:"Reaction.id"`
	Subject   []ReactionSubject `json:"Reaction.subject"`
	Creation  time.Time         `json:"Reaction.creation"`
	Author    []User            `json:"Reaction.author"`
	Message   string            `json:"Reaction.message"`
	Emotion   emotion.Emotion   `json:"Reaction.emotion"`
	Reactions []Reaction        `json:"Reaction.reactions"`
}
