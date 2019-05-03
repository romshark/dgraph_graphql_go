package dgraph

import "time"

// Session defines the Session type query object
type Session struct {
	UID      string    `json:"uid"`
	Key      string    `json:"Session.key"`
	Creation time.Time `json:"Session.creation"`
	User     []User    `json:"Session.user"`
}
