package resolver

import (
	"context"
	"demo/store"
)

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
