package http

import (
	"context"
	"log"
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
	conf         ServerConfig
	httpSrv      *http.Server
	addr         net.Addr
	onGraphQuery trn.OnGraphQuery
	onAuth       trn.OnAuth
	onDebugAuth  trn.OnDebugAuth
	onDebugSess  trn.OnDebugSess
	debugLog     *log.Logger
	errorLog     *log.Logger
}

// NewServer creates a new unencrypted JSON based HTTP transport.
// Use NewSecure to enable encryption instead
func NewServer(conf ServerConfig) (trn.Server, error) {
	if err := conf.Prepare(); err != nil {
		return nil, err
	}

	t := &Server{
		addrReadWait: &sync.WaitGroup{},
		conf:         conf,
	}
	t.httpSrv = &http.Server{
		Addr:    conf.Host,
		Handler: t,
	}

	if conf.TLS != nil {
		t.httpSrv.TLSConfig = conf.TLS.Config
	}

	t.addrReadWait.Add(1)
	return t, nil
}

// Init implements the transport.Transport interface
func (t *Server) Init(
	onGraphQuery trn.OnGraphQuery,
	onAuth trn.OnAuth,
	onDebugAuth trn.OnDebugAuth,
	onDebugSess trn.OnDebugSess,
	debugLog *log.Logger,
	errorLog *log.Logger,
) error {
	if onGraphQuery == nil {
		panic("missing onGraphQuery callback")
	}
	if onAuth == nil {
		panic("missing onAuth callback")
	}
	if onDebugAuth == nil {
		panic("missing onDebugAuth callback")
	}
	if onDebugSess == nil {
		panic("missing onDebugSess callback")
	}
	t.onGraphQuery = onGraphQuery
	t.onAuth = onAuth
	t.onDebugAuth = onDebugAuth
	t.onDebugSess = onDebugSess
	t.debugLog = debugLog
	t.errorLog = errorLog
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

	tcpListener := tcpKeepAliveListener{
		TCPListener:       listener.(*net.TCPListener),
		KeepAliveDuration: t.conf.KeepAliveDuration,
	}

	if t.conf.TLS != nil {
		t.debugLog.Print("listening https://" + t.addr.String())

		if err := t.httpSrv.ServeTLS(
			tcpListener,
			t.conf.TLS.CertificateFilePath,
			t.conf.TLS.PrivateKeyFilePath,
		); err != http.ErrServerClosed {
			return err
		}
	} else {
		t.debugLog.Print("listening http://" + t.addr.String())

		if err := t.httpSrv.Serve(tcpListener); err != http.ErrServerClosed {
			return err
		}
	}

	return nil
}

// Shutdown implements the transport.Transport interface
func (t *Server) Shutdown(ctx context.Context) error {
	return t.httpSrv.Shutdown(ctx)
}

// ServeHTTP implements the http.Handler interface
func (t *Server) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	// Authenticate the client by passing the session in the context
	// of the request
	req = t.auth(req)

	switch req.Method {
	case "POST":
		switch req.URL.Path {
		case "/g":
			t.handleGraphQuery(resp, req)
		case "/debug":
			t.handleDebugAuth(resp, req)
		default:
			// Unsupported path
			http.Error(
				resp,
				http.StatusText(http.StatusNotFound),
				http.StatusNotFound,
			)
		}
	case "GET":
		switch req.URL.Path {
		case "/playground":
			t.servePlayground(resp, req)
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

// Config returns the active configuration
func (t *Server) Config() ServerConfig {
	return ServerConfig{
		Host:              t.conf.Host,
		KeepAliveDuration: t.conf.KeepAliveDuration,
		TLS:               t.conf.TLS.Clone(),
		Playground:        t.conf.Playground,
	}
}
