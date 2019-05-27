package graph

import (
	"strconv"

	"github.com/graph-gophers/graphql-go/errors"
)

// GQLError represents a GraphQL error
type GQLError struct {
	errs []*errors.QueryError
}

func (gqlerr GQLError) Error() string {
	msg := "graph: ["
	for i, err := range gqlerr.errs {
		msg += strconv.Itoa(i) + ": " + err.Error() + "; "
	}
	msg += "]"
	return msg
}
