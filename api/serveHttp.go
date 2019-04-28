package api

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/pkg/errors"
)

type stackTracer interface {
	StackTrace() errors.StackTrace
}

// handleGraphQuery handles an HTTP graph query
func (srv *server) handleGraphQuery(
	resp http.ResponseWriter,
	req *http.Request,
) {
	// Read query body
	query, err := ioutil.ReadAll(req.Body)
	if err != nil {
		http.Error(
			resp,
			http.StatusText(http.StatusBadRequest),
			http.StatusBadRequest,
		)
	}

	// Resolve query
	reply, err := srv.graph.Query(context.Background(), query)
	if err != nil {
		http.Error(
			resp,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError,
		)

		// Retrieve error stack trace and log the error
		var tracedError string
		if tracErr, ok := err.(stackTracer); ok {
			tracedError = err.Error() + "\n"
			for _, f := range tracErr.StackTrace() {
				tracedError = fmt.Sprintf("%s%+s:%d\n", tracedError, f, f)
			}
		} else {
			tracedError = err.Error()
		}
		log.Printf("graph query: %s", tracedError)

		return
	}

	// Reply successfully
	resp.Header().Set("Content-Type", "application/json")
	if _, err := resp.Write(reply); err != nil {
		log.Printf("reply write: %s", err)
	}
}

// ServeHTTP implements the http.Handler interface
func (srv *server) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "POST":
		switch req.URL.Path {
		case "/g":
			srv.handleGraphQuery(resp, req)
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
