package api

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"github.com/romshark/dgraph_graphql_go/api/graph"
	"github.com/romshark/dgraph_graphql_go/api/graph/resolver"
	strerr "github.com/romshark/dgraph_graphql_go/store/errors"
)

type stackTracer interface {
	StackTrace() errors.StackTrace
}

func (srv *server) handleUnexpectedError(err error) {
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
	srv.logErrf("graph query: %s", tracedError)
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

		// Unexpected resolver error
		srv.handleUnexpectedError(queryErr)
		return graph.Response{}, resolverErr
	}

	if queryErr != nil {
		if gqlErr, isGQLErr := queryErr.(graph.GQLError); isGQLErr {
			return graph.Response{
				Error: &graph.ResponseError{
					Message: gqlErr.Error(),
				},
			}, nil
		}
		// Unexpected internal server error
		srv.handleUnexpectedError(queryErr)
		return graph.Response{}, queryErr
	}

	// Reply successfully
	return graph.Response{Data: replyData}, nil
}
