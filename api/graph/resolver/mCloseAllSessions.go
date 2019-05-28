package resolver

import (
	"context"

	"github.com/romshark/dgraph_graphql_go/api/graph/auth"
	"github.com/romshark/dgraph_graphql_go/store"
)

// CloseAllSessions resolves Mutation.closeAllSessions
func (rsv *Resolver) CloseAllSessions(
	ctx context.Context,
	params struct {
		User string
	},
) []string {
	if err := auth.Authorize(ctx, auth.IsOwner{
		Owner: store.ID(params.User),
	}); err != nil {
		rsv.error(ctx, err)
		return nil
	}

	result, err := rsv.str.CloseAllSessions(
		ctx,
		store.ID(params.User),
	)
	if err != nil {
		rsv.error(ctx, err)
		return nil
	}

	return result
}
