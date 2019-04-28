package api

import (
	"context"
	"demo/api/graph"
	"demo/store"
	"log"
	"net"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/pkg/errors"
)

// ServerOptions defines the API server options
type ServerOptions struct {
	Host              string
	KeepAliveDuration time.Duration
}

// SetDefaults sets the default options
func (opts *ServerOptions) SetDefaults() {
	if opts.KeepAliveDuration == time.Duration(0) {
		opts.KeepAliveDuration = 3 * time.Minute
	}
}

// Server interfaces an API server implementation
type Server interface {
	Launch() error
	Shutdown(context.Context) error
	AwaitShutdown()
	Addr() url.URL
}

type server struct {
	opts    ServerOptions
	httpSrv *http.Server
	wg      *sync.WaitGroup
	store   store.Store
	addr    net.Addr
	graph   *graph.Graph
}

// NewServer creates a new API server instance
func NewServer(opts ServerOptions, str store.Store) Server {
	opts.SetDefaults()
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
	return newSrv
}

// Run implements the Server interface
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
