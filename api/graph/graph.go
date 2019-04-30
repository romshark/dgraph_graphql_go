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
	var queryObject struct {
		Query         string                 `json:"query"`
		OperationName string                 `json:"operationName"`
		Variables     map[string]interface{} `json:"variables"`
	}
	if err := json.Unmarshal(query, &queryObject); err != nil {
		return nil, errors.New("invalid query object")
	}
	rep := graph.schema.Exec(
		ctx,
		queryObject.Query,
		queryObject.OperationName,
		queryObject.Variables,
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
