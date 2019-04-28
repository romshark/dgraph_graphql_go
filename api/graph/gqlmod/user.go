package gqlmod

import (
	"demo/store"
	"time"
)

// User defines the User type query object
type User struct {
	ID          *store.ID  `json:"id"`
	Creation    *time.Time `json:"creation"`
	Email       *string    `json:"email"`
	DisplayName *string    `json:"displayName"`
	Posts       []Post     `json:"posts"`
}
