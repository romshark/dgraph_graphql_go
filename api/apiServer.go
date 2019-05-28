package api

import (
	"context"
	"fmt"
	"sync"

	"github.com/pkg/errors"
	"github.com/romshark/dgraph_graphql_go/api/graph"
	"github.com/romshark/dgraph_graphql_go/api/options"
	"github.com/romshark/dgraph_graphql_go/api/transport"
	"github.com/romshark/dgraph_graphql_go/api/validator"
	"github.com/romshark/dgraph_graphql_go/store"
	"github.com/romshark/dgraph_graphql_go/store/dgraph"
)

// Server interfaces an API server implementation
type Server interface {
	// Launche starts the API server and unblocks as soon as the server
	// is running
	Launch() error

	// Shutdown instructs the server to shut down gracefully and blocks until
	// the server is shut down
	Shutdown(context.Context) error

	// AwaitShutdown blocks until the server is shut down
	AwaitShutdown()
}

type server struct {
	opts                 options.ServerOptions
	store                store.Store
	graph                *graph.Graph
	debugSessionKey      []byte
	transports           []transport.Server
	shutdownAwaitBlocker *sync.WaitGroup
}

// NewServer creates a new API server instance
func NewServer(opts options.ServerOptions) (Server, error) {
	if err := opts.Prepare(); err != nil {
		return nil, fmt.Errorf("options: %s", err)
	}

	// Initialize validator
	validator, err := validator.NewValidator(
		opts.Mode == options.ModeProduction,
		validator.Options{
			PasswordLenMin:        6,
			PasswordLenMax:        256,
			EmailLenMax:           96,
			PostContentsLenMin:    1,
			PostContentsLenMax:    256,
			PostTitleLenMin:       2,
			PostTitleLenMax:       64,
			ReactionMessageLenMin: 1,
			ReactionMessageLenMax: 256,
			UserDisplayNameLenMin: 2,
			UserDisplayNameLenMax: 64,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("validator init: %s", err)
	}

	// Initialize store instance
	store := dgraph.NewStore(
		opts.DBHost,

		// Compare password
		func(hash, password string) bool {
			return opts.PasswordHasher.Compare([]byte(hash), []byte(password))
		},

		opts.DebugLog,
		opts.ErrorLog,
	)

	graph, err := graph.New(
		store,
		validator,
		opts.SessionKeyGenerator,
		opts.PasswordHasher,
	)
	if err != nil {
		return nil, errors.Wrap(err, "graph init")
	}

	// Initialize API server instance
	newSrv := &server{
		store:                store,
		opts:                 opts,
		graph:                graph,
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
			opts.DebugLog,
			opts.ErrorLog,
		); err != nil {
			return nil, err
		}
	}
	opts.DebugLog.Print("all transports initialized")

	// Generate the debug user session key if the debug user is enabled
	if opts.DebugUser.Status != options.DebugUserDisabled {
		newSrv.debugSessionKey = []byte(opts.SessionKeyGenerator.Generate())
	}

	return newSrv, nil
}

func (srv *server) logErrf(format string, v ...interface{}) {
	srv.opts.ErrorLog.Printf(format, v...)
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
				srv.logErrf("transport: %s", err)
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
				srv.logErrf("transport shutdown: %s", err)
			}
			wg.Done()
		}()
	}
	if len(shutdownErrs) < 1 {
		return nil
	}
	return errors.Errorf("ERR: shutdown: %v", shutdownErrs)
}
