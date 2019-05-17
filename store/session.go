package store

import "time"

// Session represents a Session entity
type Session struct {
	GraphNode

	Key      string
	Creation time.Time
	User     *User
}
