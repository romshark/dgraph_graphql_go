package resolver

import (
	"context"

	"github.com/romshark/dgraph_graphql_go/store"
)

// CtxKey represents a context.Context value key type
type CtxKey int

// CtxErrorRef defines the context.Context error reference value key
const CtxErrorRef CtxKey = 1

// Resolver represents the root Graph resolver
type Resolver struct {
	str store.Store
}

// New creates a new graph resolver instance
func New(str store.Store) *Resolver {
	return &Resolver{
		str: str,
	}
}

// Users resolves Query.users
func (rsv *Resolver) Users(ctx context.Context) ([]*User, error) {
	return nil, nil
}

// Posts resolves Query.posts
func (rsv *Resolver) Posts(ctx context.Context) ([]*Post, error) {
	return nil, nil
}

// error writes an error to the resolver context for the API server to read
func (rsv *Resolver) error(ctx context.Context, err error) {
	ctxErr := ctx.Value(CtxErrorRef).(*error)
	*ctxErr = err
}
