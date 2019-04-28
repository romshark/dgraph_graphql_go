package graph

import (
	"context"
	"encoding/json"

	"github.com/graph-gophers/graphql-go"
	"github.com/pkg/errors"
	rsv "github.com/romshark/dgraph_graphql_go/api/graph/resolver"
	"github.com/romshark/dgraph_graphql_go/store"
)

// Graph represents the graph resolution engine
type Graph struct {
	resolver *rsv.Resolver
	schema   *graphql.Schema
}

// New creates a new graph resolver instance
func New(str store.Store) *Graph {
	rsv := rsv.New(str)
	shm := graphql.MustParseSchema(schema, rsv)
	return &Graph{
		resolver: rsv,
		schema:   shm,
	}
}

// QueryError represents a graph query error
type QueryError struct {
	Code    string `json:"c"`
	Message string `json:"m"`
}

// Query executes a graph query and returns a JSON result (or error)
func (graph *Graph) Query(
	ctx context.Context,
	query []byte,
) (reply []byte, err error) {
	var params struct {
		Query         string                 `json:"query"`
		OperationName string                 `json:"operationName"`
		Variables     map[string]interface{} `json:"variables"`
	}
	if err := json.Unmarshal(query, &params); err != nil {
		return []byte(""), errors.Wrap(err, "unmarshalling query")
	}
	rep := graph.schema.Exec(
		ctx,
		params.Query,
		params.OperationName,
		params.Variables,
	)

	// Serialize errors
	if rep.Errors != nil {
		errs := make([]QueryError, len(rep.Errors))
		for index, err := range rep.Errors {
			typedErr, isExpectedErr := store.ParseError(err.Message)
			if !isExpectedErr {
				return nil, errors.Errorf("unexpected error: %s", err.Message)
			}
			errs[index] = QueryError{
				Code:    typedErr.Code,
				Message: typedErr.Message,
			}
		}

		jsonErr, err := json.Marshal(errs)
		if err != nil {
			return nil, errors.Wrap(err, "marshal graph query errors")
		}

		return jsonErr, nil
	}

	reply, err = json.Marshal(rep)
	if err != nil {
		return []byte(""), errors.Wrap(err, "failed marshalling JSON response")
	}

	return reply, nil
}
