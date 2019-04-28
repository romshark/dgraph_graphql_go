package gqlmod

import (
	"demo/store"
	"demo/store/enum/emotion"
	"time"
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
