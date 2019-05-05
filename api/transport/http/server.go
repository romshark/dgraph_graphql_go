package http

import (
	"context"
	"net"
	"net/http"
	"net/url"
	"sync"

	"github.com/pkg/errors"
	trn "github.com/romshark/dgraph_graphql_go/api/transport"
)

// Server represents an HTTP based server transport implementation
type Server struct {
	addrReadWait *sync.WaitGroup
	opts         ServerOptions
	httpSrv      *http.Server
	addr         net.Addr
	onGraphQuery trn.OnGraphQuery
	onRootAuth   trn.OnRootAuth
}

// NewServer creates a new unencrypted JSON based HTTP transport.
// Use NewSecure to enable encryption instead
func NewServer(opts ServerOptions) trn.Server {
	t := &Server{
		addrReadWait: &sync.WaitGroup{},
		opts:         opts,
	}
	t.httpSrv = &http.Server{
		Addr:    opts.Host,
		Handler: t,
	}
	t.addrReadWait.Add(1)
	return t
}

// NewSecureServer creates a new TLS protected JSON based HTTP transport
func NewSecureServer(opts ServerOptions) trn.Server {
	// TODO: add TLS support
	t := NewServer(opts)
	return t
}

// Init implements the transport.Transport interface
func (t *Server) Init(
	onGraphQuery trn.OnGraphQuery,
	onRootAuth trn.OnRootAuth,
) error {
	if onGraphQuery == nil {
		panic("missing onGraphQuery callback")
	}
	if onRootAuth == nil {
		panic("missing onRootAuth callback")
	}
	t.onGraphQuery = onGraphQuery
	t.onRootAuth = onRootAuth
	return nil
}

// Run implements the transport.Transport interface
func (t *Server) Run() error {
	addr := t.httpSrv.Addr
	if addr == "" {
		addr = ":http"
	}
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return errors.Wrap(err, "TCP listener setup")
	}

	t.addr = listener.Addr()
	// Address determined, readers must be unblocked
	t.addrReadWait.Done()

	if err := t.httpSrv.Serve(tcpKeepAliveListener{
		TCPListener:       listener.(*net.TCPListener),
		KeepAliveDuration: t.opts.KeepAliveDuration,
	}); err != http.ErrServerClosed {
		return err
	}
	return nil
}

// Shutdown implements the transport.Transport interface
func (t *Server) Shutdown(ctx context.Context) error {
	return t.httpSrv.Shutdown(ctx)
}

// ServeHTTP implements the http.Handler interface
func (t *Server) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "POST":
		switch req.URL.Path {
		case "/g":
			t.handleGraphQuery(resp, req)
		case "/root":
			t.handleRootAuth(resp, req)
		default:
			// Unsupported path
			http.Error(
				resp,
				http.StatusText(http.StatusNotFound),
				http.StatusNotFound,
			)
		}
	default:
		http.Error(resp, "unsupported method", http.StatusMethodNotAllowed)
	}
}

// Addr returns the host address URL.
// Blocks until the listener is initialized and the actual address is known
func (t *Server) Addr() url.URL {
	t.addrReadWait.Wait()
	hostAddrStr := t.addr.String()
	return url.URL{
		Scheme: "http",
		Host:   hostAddrStr,
		Path:   "/",
	}
}
