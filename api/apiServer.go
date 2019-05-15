package api

import (
	"context"
	"log"
	"sync"

	"github.com/pkg/errors"
	"github.com/romshark/dgraph_graphql_go/api/graph"
	"github.com/romshark/dgraph_graphql_go/api/transport"
	"github.com/romshark/dgraph_graphql_go/store"
	"github.com/romshark/dgraph_graphql_go/store/dgraph"
)

// Server interfaces an API server implementation
type Server interface {
	// Launche starts the API server and unblocks as soon as the server
	// is running
	Launch() error

	// Shutdown instructs the server to shut down gracefuly and blocks until
	// the server is shut down
	Shutdown(context.Context) error

	// AwaitShutdown blocks until the server is shut down
	AwaitShutdown()
}

type server struct {
	opts                 ServerOptions
	store                store.Store
	graph                *graph.Graph
	debugSessionKey      []byte
	transports           []transport.Server
	shutdownAwaitBlocker *sync.WaitGroup
}

// NewServer creates a new API server instance
func NewServer(opts ServerOptions) (Server, error) {
	opts.SetDefaults()

	// Initialize store instance
	str := dgraph.NewStore(
		opts.DBHost,
		opts.SessionKeyGenerator,
		opts.PasswordHasher,
	)

	// Initialize API server instance
	newSrv := &server{
		store:                str,
		opts:                 opts,
		graph:                graph.New(str),
		transports:           opts.Transport,
		shutdownAwaitBlocker: &sync.WaitGroup{},
	}

	// Initialize transports
	for _, transport := range opts.Transport {
		if err := transport.Init(
			newSrv.onGraphQuery,
			newSrv.onAuth,
			newSrv.onDebugAuth,
			newSrv.onDebugSess,
		); err != nil {
			return nil, err
		}
	}

	// Generate the debug user session key if the debug user is enabled
	if opts.DebugUser.Status != DebugUserDisabled {
		newSrv.debugSessionKey = []byte(opts.SessionKeyGenerator.Generate())
	}

	return newSrv, nil
}

// Launch implements the Server interface
func (srv *server) Launch() error {
	// Prepare the store
	if err := srv.store.Prepare(); err != nil {
		return errors.Wrap(err, "store preparation")
	}

	// Launch all transports
	srv.shutdownAwaitBlocker.Add(len(srv.transports))
	for _, transport := range srv.transports {
		t := transport
		go func() {
			if err := t.Run(); err != nil {
				log.Printf("ERR: transport: %s", err)
			}
			srv.shutdownAwaitBlocker.Done()
		}()
	}

	return nil
}

// AwaitShutdown implements the Server interface
func (srv *server) AwaitShutdown() {
	srv.shutdownAwaitBlocker.Wait()
}

// Shutdown implements the Server interface
func (srv *server) Shutdown(ctx context.Context) error {
	wg := &sync.WaitGroup{}
	wg.Add(len(srv.transports))

	shutdownErrs := make([]error, 0)
	for _, transport := range srv.transports {
		t := transport
		go func() {
			if err := t.Shutdown(ctx); err != nil {
				shutdownErrs = append(
					shutdownErrs,
					errors.Wrap(err, "transport shutdown"),
				)
				log.Printf("ERR: transport shutdown: %s", err)
			}
			wg.Done()
		}()
	}
	if len(shutdownErrs) < 1 {
		return nil
	}
	return errors.Errorf("ERR: shutdown: %v", shutdownErrs)
}
