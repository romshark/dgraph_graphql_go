package resolver

import (
	"context"

	"github.com/romshark/dgraph_graphql_go/api/graph/auth"
)

// CloseSession resolves Mutation.closeSession
func (rsv *Resolver) CloseSession(
	ctx context.Context,
	params struct {
		Key string
	},
) (bool, error) {
	if err := auth.Authorize(ctx, auth.IsUser{}); err != nil {
		rsv.error(ctx, err)
		return false, err
	}

	result, err := rsv.str.CloseSession(
		ctx,
		params.Key,
	)
	if err != nil {
		rsv.error(ctx, err)
		return false, err
	}

	return result, nil
}
