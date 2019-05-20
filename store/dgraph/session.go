package dgraph

import "time"

// Session represents a database model for the Session entity
type Session struct {
	UID       string    `json:"uid"`
	Key       string    `json:"Session.key"`
	Creation  time.Time `json:"Session.creation"`
	User      []User    `json:"Session.user"`
	RSessions []UID     `json:"~sessions"`
}
