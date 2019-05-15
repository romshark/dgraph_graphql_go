package gqlmod

import (
	"time"
)

// Session defines the Session type query object
type Session struct {
	Key      *string    `json:"key"`
	Creation *time.Time `json:"creation"`
	User     *User      `json:"user"`
}
