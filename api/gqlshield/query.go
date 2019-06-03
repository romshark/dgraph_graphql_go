package gqlshield

import (
	"sort"
	"time"
)

// Query represents a query object
type Query interface {
	// Query returns a copy of the query string
	Query() []byte

	// ID returns the unique identifier of the query
	ID() ID

	// Creation returns the time of creation
	Creation() time.Time

	// Name returns the query name
	Name() string

	// Parameters returns a copy of the list of query parameters
	Parameters() map[string]Parameter

	// WhitelistedFor returns the copy of role IDs the query is whitelisted for
	WhitelistedFor() []int
}

// Parameter represents a query parameter
type Parameter struct {
	MaxValueLength uint32 `json:"max-value-length"`
}

// query represents a whitelisted query
type query struct {
	id             ID
	query          []byte
	creation       time.Time
	name           string
	parameters     map[string]Parameter
	whitelistedFor map[int]struct{}
}

func (q *query) ID() ID {
	return q.id
}

func (q *query) Creation() time.Time {
	return q.creation
}

func (q *query) Name() string {
	return q.name
}

func (q *query) Query() []byte {
	query := make([]byte, len(q.query))
	copy(query, q.query)
	return query
}

func (q *query) Parameters() map[string]Parameter {
	params := make(map[string]Parameter, len(q.parameters))
	for name, param := range q.parameters {
		params[name] = param
	}
	return params
}

func (q *query) WhitelistedFor() []int {
	clientRoleIDs := make([]int, len(q.whitelistedFor))
	index := 0
	for roleID := range q.whitelistedFor {
		clientRoleIDs[index] = roleID
		index++
	}
	sort.Ints(clientRoleIDs)
	return clientRoleIDs
}
