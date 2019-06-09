package api

import (
	"context"
	"fmt"
	"sync"

	"github.com/romshark/dgraph_graphql_go/api/graph/auth"

	"github.com/pkg/errors"
	"github.com/romshark/dgraph_graphql_go/api/config"
	"github.com/romshark/dgraph_graphql_go/api/gqlshield"
	"github.com/romshark/dgraph_graphql_go/api/graph"
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
	conf                 *config.ServerConfig
	store                store.Store
	graph                *graph.Graph
	debugSessionKey      []byte
	transports           []transport.Server
	shutdownAwaitBlocker *sync.WaitGroup
}

// NewServer creates a new API server instance
func NewServer(conf *config.ServerConfig) (Server, error) {
	if err := conf.Prepare(); err != nil {
		return nil, fmt.Errorf("config: %s", err)
	}

	// Initialize validator
	validator, err := validator.NewValidator(
		conf.Mode == config.ModeProduction,
		validator.Config{
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
		conf.DBHost,

		// Compare password
		func(hash, password string) bool {
			return conf.PasswordHasher.Compare([]byte(hash), []byte(password))
		},

		conf.DebugLog,
		conf.ErrorLog,
	)

	// Initialize the GraphQL shield persistency manager
	var shieldPersistencyManager gqlshield.PersistencyManager
	if conf.Shield.PersistencyFilePath != "" {
		manager, err := gqlshield.NewPepersistencyManagerFileJSON(
			conf.Shield.PersistencyFilePath,
			true,
		)
		if err != nil {
			return nil, errors.Wrap(err, "GraphQL shield persistency manager init")
		}
		shieldPersistencyManager = manager
	}

	queryWhitelistingEnabled := gqlshield.WhitelistDisabled
	if conf.Shield.WhitelistEnabled {
		queryWhitelistingEnabled = gqlshield.WhitelistEnabled
	}

	graphShield, err := gqlshield.NewGraphQLShield(
		gqlshield.Config{
			WhitelistOption:    queryWhitelistingEnabled,
			PersistencyManager: shieldPersistencyManager,
		},
		gqlshield.ClientRole{
			ID:   int(auth.GQLShieldClientDebug),
			Name: "debug",
		},
		gqlshield.ClientRole{
			ID:   int(auth.GQLShieldClientGuest),
			Name: "guest",
		},
		gqlshield.ClientRole{
			ID:   int(auth.GQLShieldClientRegular),
			Name: "regular",
		},
	)
	if err != nil {
		return nil, errors.Wrap(err, "graph shield init")
	}

	graph, err := graph.New(
		store,
		validator,
		conf.SessionKeyGenerator,
		conf.PasswordHasher,
		graphShield,
	)
	if err != nil {
		return nil, errors.Wrap(err, "graph init")
	}

	// Initialize API server instance
	newSrv := &server{
		store:                store,
		conf:                 conf,
		graph:                graph,
		transports:           conf.Transport,
		shutdownAwaitBlocker: &sync.WaitGroup{},
	}

	// Initialize transports
	for _, transport := range conf.Transport {
		if err := transport.Init(
			newSrv.onGraphQuery,
			newSrv.onAuth,
			newSrv.onDebugAuth,
			newSrv.onDebugSess,
			conf.DebugLog,
			conf.ErrorLog,
		); err != nil {
			return nil, err
		}
	}
	conf.DebugLog.Print("all transports initialized")

	// Generate the debug user session key if the debug user is enabled
	if conf.DebugUser.Mode != config.DebugUserDisabled {
		newSrv.debugSessionKey = []byte(conf.SessionKeyGenerator.Generate())
	}

	return newSrv, nil
}

func (srv *server) logErrf(format string, v ...interface{}) {
	srv.conf.ErrorLog.Printf(format, v...)
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
