package api

import (
	"context"
	"fmt"
	"log"

	"github.com/pkg/errors"
	"github.com/romshark/dgraph_graphql_go/api/graph"
	"github.com/romshark/dgraph_graphql_go/api/graph/resolver"
	strerr "github.com/romshark/dgraph_graphql_go/store/errors"
)

type stackTracer interface {
	StackTrace() errors.StackTrace
}

// onGraphQuery handles a graph query
func (srv *server) onGraphQuery(
	ctx context.Context,
	query graph.Query,
) (graph.Response, error) {
	// Resolve query
	var resolverErr error
	ctxWithRsvErr := context.WithValue(
		ctx,
		resolver.CtxErrorRef,
		&resolverErr,
	)
	replyData, queryErr := srv.graph.Query(ctxWithRsvErr, query)

	errCode := strerr.ErrorCode(resolverErr)
	if resolverErr != nil {
		if errCode != "" {
			// Expected error
			return graph.Response{
				Error: &graph.ResponseError{
					Code:    errCode,
					Message: resolverErr.Error(),
				},
			}, nil
		}

		// Unexpected error

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

		return graph.Response{}, resolverErr
	}
	if queryErr != nil {
		return graph.Response{}, queryErr
	}

	// Reply successfully
	return graph.Response{Data: replyData}, nil
}
