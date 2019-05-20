package transport

import (
	"github.com/romshark/dgraph_graphql_go/api/graph/gqlmod"
)

// Client defines the interface of the transport layer implementation.
// Run and Init are not intended to be thread-safe and shall only be used
// by a single goroutine!
type Client interface {
	// SignIn signs the client into a user
	SignIn(email, password string) (*gqlmod.Session, error)

	// SignInDebug signs the client into a debug user
	SignInDebug(username, password string) error

	// Auth authenticates the client by a user session key
	Auth(sessionKey string) (*gqlmod.Session, error)

	// Query performs an API query
	Query(query string, result interface{}) error

	// QueryVar performs a parameterized API query
	QueryVar(
		query string,
		vars map[string]interface{},
		result interface{},
	) error
}
