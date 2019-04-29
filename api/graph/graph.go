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

// Query executes a graph query and returns a JSON encoded result (or an error)
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
		return nil, errors.Wrap(err, "unmarshalling query")
	}
	rep := graph.schema.Exec(
		ctx,
		params.Query,
		params.OperationName,
		params.Variables,
	)

	if rep.Errors != nil {
		// Serialize errors
		errMsg := "graph error:"
		for _, err := range rep.Errors {
			errMsg += " " + err.Error() + ";"
		}
		return nil, errors.New(errMsg)
	}

	return rep.Data, nil
}
