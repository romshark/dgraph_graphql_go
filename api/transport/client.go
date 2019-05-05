package transport

// Client defines the interface of the transport layer implementation.
// Run and Init are not intended to be thread-safe and shall only be used
// by a single goroutine!
type Client interface {
	// Authenticates as root user
	AuthRoot(username, password string) error

	// Query performs an API query
	Query(query string, result interface{}) error

	// QueryVar performs a parameterized API query
	QueryVar(query string, vars map[string]string, result interface{}) error
}
