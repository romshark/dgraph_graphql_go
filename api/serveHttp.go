package api

import (
	"net/http"

	"github.com/pkg/errors"
)

var dataRepHead = []byte("{\"d\":")
var dataRepTail = []byte("}")

type stackTracer interface {
	StackTrace() errors.StackTrace
}

// responseError represents a response error object
type responseError struct {
	Code    string `json:"c,omitempty"`
	Message string `json:"m,omitempty"`
}

// response represents a response object
type response struct {
	Data  string         `json:"d,omitempty"`
	Error *responseError `json:"e,omitempty"`
}

// ResponseError represents a response error object
type ResponseError struct {
	Code    string `json:"c,omitempty"`
	Message string `json:"m,omitempty"`
}

// Response represents a response object
type Response struct {
	Data  interface{}    `json:"d,omitempty"`
	Error *ResponseError `json:"e,omitempty"`
}

// ServeHTTP implements the http.Handler interface
func (srv *server) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "POST":
		switch req.URL.Path {
		case "/g":
			srv.handleGraphQuery(resp, req)
		case "/root":
			srv.handleRoot(resp, req)
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
