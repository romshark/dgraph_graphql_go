package gqlshield

import "time"

// QueryModel represents the model of a query entry for serialization
type QueryModel struct {
	Query          string               `json:"query"`
	Creation       time.Time            `json:"creation"`
	Name           string               `json:"name"`
	Parameters     map[string]Parameter `json:"parameters"`
	WhitelistedFor []int                `json:"whitelisted-for"`
}

// State represents the state of the GraphQL shield
type State struct {
	Roles              []ClientRole          `json:"roles"`
	WhitelistedQueries map[string]QueryModel `json:"whitelisted-queries"`
}
