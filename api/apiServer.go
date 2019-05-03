package api

import (
	"context"
	"log"
	"net"
	"net/http"
	"net/url"
	"sync"

	"github.com/pkg/errors"
	"github.com/romshark/dgraph_graphql_go/api/graph"
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

	// Addr returns the address the API server is serving on
	Addr() url.URL
}

type server struct {
	opts           ServerOptions
	httpSrv        *http.Server
	wg             *sync.WaitGroup
	store          store.Store
	addr           net.Addr
	graph          *graph.Graph
	rootSessionKey []byte
}

// NewServer creates a new API server instance
func NewServer(opts ServerOptions) Server {
	opts.SetDefaults()

	// Initialize store instance
	str := dgraph.NewStore(
		opts.DBHost,
		opts.SessionKeyGenerator,
		opts.PasswordHasher,
	)

	// Initialize API server instance
	newSrv := &server{
		store: str,
		opts:  opts,
		wg:    &sync.WaitGroup{},
		graph: graph.New(str),
	}
	newSrv.wg.Add(1)
	newSrv.httpSrv = &http.Server{
		Addr:    opts.Host,
		Handler: newSrv,
	}

	// Generate the root session key if the root user is enabled
	if opts.RootUser.Status != RootUserDisabled {
		newSrv.rootSessionKey = []byte(opts.SessionKeyGenerator.Generate())
	}

	return newSrv
}

// Launch implements the Server interface
func (srv *server) Launch() error {
	if err := srv.store.Prepare(); err != nil {
		return errors.Wrap(err, "store preparation")
	}

	addr := srv.httpSrv.Addr
	if addr == "" {
		addr = ":http"
	}
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return errors.Wrap(err, "TCP listener setup")
	}
	srv.addr = listener.Addr()
	go func() {
		if err := srv.httpSrv.Serve(tcpKeepAliveListener{
			TCPListener:       listener.(*net.TCPListener),
			KeepAliveDuration: srv.opts.KeepAliveDuration,
		}); err != http.ErrServerClosed {
			log.Fatalf("http serve: %s", err)
		}
		srv.wg.Done()
	}()
	return nil
}

// AwaitShutdown implements the Server interface
func (srv *server) AwaitShutdown() {
	srv.wg.Wait()
}

// Shutdown implements the Server interface
func (srv *server) Shutdown(ctx context.Context) error {
	return srv.httpSrv.Shutdown(ctx)
}

// Addr implements the Server interface
func (srv *server) Addr() url.URL {
	return url.URL{
		Scheme: "http",
		Host:   srv.addr.String(),
		Path:   "/",
	}
}
