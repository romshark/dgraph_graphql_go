package store

import "time"

// Post represents a Post entity
type Post struct {
	GraphNode

	ID        ID
	Creation  time.Time
	Author    *User
	Title     string
	Contents  string
	Reactions []Reaction
}
