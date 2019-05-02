package gqlmod

import (
	"time"

	"github.com/romshark/dgraph_graphql_go/store"
)

// Session defines the Session type query object
type Session struct {
	ID       *store.ID  `json:"id"`
	Key      *string    `json:"key"`
	Creation *time.Time `json:"creation"`
	User     *User      `json:"user"`
}
