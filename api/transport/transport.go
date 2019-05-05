package transport

import (
	"context"

	"github.com/romshark/dgraph_graphql_go/api/graph"
)

// OnGraphQuery defines the graph query callback function
type OnGraphQuery func(context.Context, graph.Query) (graph.Response, error)

// OnRootAuth defines the root authentication callback function
type OnRootAuth func(username, password string) ([]byte, bool)

// Server defines the interface of the server transport layer implementation.
// Run and Init are not intended to be thread-safe and shall only be used
// by a single goroutine!
type Server interface {
	// Init initializes the server transport implementation.
	// The provided callbacks must be registered and invoked accordingly
	Init(onGraphQuery OnGraphQuery, onRootAuth OnRootAuth) error

	// Run starts serving. Blocks until the underlying server is shut down
	Run() error

	// Shutdown instructs the underlying server to shut down and blocks until
	// it is
	Shutdown(ctx context.Context) error
}
