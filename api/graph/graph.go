package graph

import (
	"context"
	"errors"
	"fmt"

	"github.com/graph-gophers/graphql-go"
	"github.com/romshark/dgraph_graphql_go/api/graph/resolver"
	rsv "github.com/romshark/dgraph_graphql_go/api/graph/resolver"
	"github.com/romshark/dgraph_graphql_go/api/passhash"
	"github.com/romshark/dgraph_graphql_go/api/sesskeygen"
	"github.com/romshark/dgraph_graphql_go/api/validator"
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
) ([]byte, error) {
	// Validate query
	if errs := graph.schema.Validate(query.Query); errs != nil {
		return nil, GQLError{errs: errs}
	}

	var resolverErr error

	// Execute query
	rep := graph.schema.Exec(
		context.WithValue(
			ctx,
			resolver.CtxErrorRef,
			&resolverErr,
		),
		query.Query,
		query.OperationName,
		query.Variables,
	)

	if resolverErr != nil {
		return nil, resolverErr
	}

	if rep.Errors != nil {
		err := GQLError{errs: rep.Errors}
		// Return an untyped error to prevent it from leaking to the client
		// because GQLError instances are served as error feedback to clients
		// but unexpected errors should not be leaked!
		return nil, errors.New(err.Error())
	}

	return rep.Data, nil
}
