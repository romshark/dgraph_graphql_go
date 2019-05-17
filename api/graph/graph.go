package graph

import (
	"context"
	"fmt"

	"github.com/romshark/dgraph_graphql_go/api/graph/validator"

	"github.com/graph-gophers/graphql-go"
	"github.com/pkg/errors"
	rsv "github.com/romshark/dgraph_graphql_go/api/graph/resolver"
	"github.com/romshark/dgraph_graphql_go/api/passhash"
	"github.com/romshark/dgraph_graphql_go/api/sesskeygen"
	"github.com/romshark/dgraph_graphql_go/store"
)

// Graph represents the graph resolution engine
type Graph struct {
	resolver *rsv.Resolver
	schema   *graphql.Schema
}

// Query represents the graph query structure
type Query struct {
	Query         string
	OperationName string
	Variables     map[string]interface{}
}

// ResponseError represents a response error object
type ResponseError struct {
	Code    string
	Message string
}

// Error implements the error interface
func (err *ResponseError) Error() string {
	return fmt.Sprintf("%s: %s", err.Code, err.Message)
}

// Response represents a response object
type Response struct {
	Data  []byte
	Error *ResponseError
}

// New creates a new graph resolver instance
func New(
	str store.Store,
	validator validator.Validator,
	sessionKeyGenerator sesskeygen.SessionKeyGenerator,
	passwordHasher passhash.PasswordHasher,
) (*Graph, error) {
	rsv, err := rsv.New(
		str,
		validator,
		sessionKeyGenerator,
		passwordHasher,
	)
	if err != nil {
		return nil, err
	}
	shm := graphql.MustParseSchema(schema, rsv)
	return &Graph{
		resolver: rsv,
		schema:   shm,
	}, nil
}

// Query executes a graph query and returns a JSON encoded result (or an error)
func (graph *Graph) Query(
	ctx context.Context,
	query Query,
) (reply []byte, err error) {
	rep := graph.schema.Exec(
		ctx,
		query.Query,
		query.OperationName,
		query.Variables,
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
