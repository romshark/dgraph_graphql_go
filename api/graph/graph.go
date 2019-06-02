package graph

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/romshark/dgraph_graphql_go/api/graph/auth"

	"github.com/graph-gophers/graphql-go"
	"github.com/romshark/dgraph_graphql_go/api/gqlshield"
	"github.com/romshark/dgraph_graphql_go/api/graph/resolver"
	rsv "github.com/romshark/dgraph_graphql_go/api/graph/resolver"
	"github.com/romshark/dgraph_graphql_go/api/passhash"
	"github.com/romshark/dgraph_graphql_go/api/sesskeygen"
	"github.com/romshark/dgraph_graphql_go/api/validator"
	"github.com/romshark/dgraph_graphql_go/store"
	strerr "github.com/romshark/dgraph_graphql_go/store/errors"
)

// Graph represents the graph resolution engine
type Graph struct {
	resolver *rsv.Resolver
	schema   *graphql.Schema
	shield   gqlshield.GraphQLShield
}

// Query represents the graph query structure
type Query struct {
	Query         json.RawMessage
	OperationName string
	Variables     map[string]string
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
	shield gqlshield.GraphQLShield,
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
		shield:   shield,
	}, nil
}

// Query executes a graph query and returns a JSON encoded result (or an error)
func (graph *Graph) Query(
	ctx context.Context,
	query Query,
) ([]byte, error) {
	clientRole := auth.GQLShieldClientRegular

	// Try to read the shield client role identifier from session
	if session, isSession := ctx.Value(
		auth.CtxSession,
	).(*auth.RequestSession); isSession {
		clientRole = session.ShieldClientRole
	}

	// Ensure the query is whitelisted for this client
	// and the arguments are valid
	if err := graph.shield.Check(
		int(clientRole),
		[]byte(query.Query),
		query.Variables,
	); err != nil {
		switch gqlshield.ErrCode(err) {
		case gqlshield.ErrWrongInput:
			return nil, strerr.New(
				strerr.ErrUnauthorized,
				err.Error(),
			)
		case gqlshield.ErrUnauthorized:
			return nil, strerr.New(
				strerr.ErrUnauthorized,
				"query blacklisted for this client",
			)
		default:
			// Unexpected error
			return nil, err
		}
	}

	args := make(map[string]interface{}, len(query.Variables))
	for name, val := range query.Variables {
		args[name] = val
	}

	// Validate query
	queryStr := string(query.Query)
	if errs := graph.schema.Validate(queryStr); errs != nil {
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
		queryStr,
		query.OperationName,
		args,
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
