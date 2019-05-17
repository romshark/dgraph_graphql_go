package store

import (
	"time"

	"github.com/romshark/dgraph_graphql_go/store/enum/emotion"
)

// Reaction represents a Reaction entity
type Reaction struct {
	GraphNode

	ID        ID
	Subject   AGraphNode
	Creation  time.Time
	Author    *User
	Message   string
	Emotion   emotion.Emotion
	Reactions []Reaction
}
