package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/pkg/errors"
	"github.com/romshark/dgraph_graphql_go/api/graph/resolver"
	strerr "github.com/romshark/dgraph_graphql_go/store/errors"
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

	jsonEncoder := json.NewEncoder(resp)

	// Resolve query
	var resolverErr error
	ctxWithRsvErr := context.WithValue(
		req.Context(),
		resolver.CtxErrorRef,
		&resolverErr,
	)
	replyData, queryErr := srv.graph.Query(ctxWithRsvErr, query)

	errCode := strerr.ErrorCode(resolverErr)
	if resolverErr != nil {
		if errCode != "" {
			// User error
			resp.WriteHeader(http.StatusBadRequest)

			respErr := responseError{
				Code:    errCode,
				Message: resolverErr.Error(),
			}
			if err := jsonEncoder.Encode(response{Error: &respErr}); err != nil {
				panic(fmt.Errorf("response JSON encode: %s", err))
			}
			return
		}

		// Unexpected internal error
		http.Error(
			resp,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError,
		)

		// Retrieve error stack trace and log the error
		var tracedError string
		if tracErr, ok := resolverErr.(stackTracer); ok {
			tracedError = resolverErr.Error() + "\n"
			for _, f := range tracErr.StackTrace() {
				tracedError = fmt.Sprintf("%s%+s:%d\n", tracedError, f, f)
			}
		} else {
			tracedError = resolverErr.Error()
		}
		log.Printf("graph query: %s", tracedError)
		return
	}

	if queryErr != nil {
		// User error
		resp.WriteHeader(http.StatusBadRequest)

		respErr := responseError{
			Message: queryErr.Error(),
		}
		if err := jsonEncoder.Encode(response{Error: &respErr}); err != nil {
			panic(fmt.Errorf("response JSON encode: %s", err))
		}
		return
	}

	// Reply successfully
	resp.Header().Set("Content-Type", "application/json")

	if _, err := resp.Write(dataRepHead); err != nil {
		log.Printf("reply data write head: %s", err)
		return
	}
	if _, err := resp.Write(replyData); err != nil {
		log.Printf("reply data write: %s", err)
		return
	}
	if _, err := resp.Write(dataRepTail); err != nil {
		log.Printf("reply data write tail: %s", err)
		return
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
