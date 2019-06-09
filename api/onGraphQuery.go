package api

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"github.com/romshark/dgraph_graphql_go/api/graph"
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
	replyData, err := srv.graph.Query(ctx, query)

	if err != nil {
		errCode := strerr.ErrorCode(err)
		if errCode != "" {
			// Expected user error
			return graph.Response{
				Error: &graph.ResponseError{
					Code:    errCode,
					Message: err.Error(),
				},
			}, nil
		}

		if gqlErr, isGQLErr := err.(graph.GQLError); isGQLErr {
			// Expected GraphQL error
			return graph.Response{
				Error: &graph.ResponseError{
					Message: gqlErr.Error(),
				},
			}, nil
		}

		// Unexpected internal server error
		srv.handleUnexpectedError(err)
		return graph.Response{}, err
	}

	// Reply successfully
	return graph.Response{Data: replyData}, nil
}
