package dgraph

import (
	"context"
	"encoding/json"

	"github.com/dgraph-io/dgo"
	"github.com/dgraph-io/dgo/protos/api"
	"github.com/pkg/errors"
	strerr "github.com/romshark/dgraph_graphql_go/store/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type transaction interface {
	Query(
		ctx context.Context,
		query string,
		res interface{},
	) error

	QueryVars(
		ctx context.Context,
		query string,
		vars map[string]string,
		res interface{},
	) error

	Mutation(
		ctx context.Context,
		mutation *api.Mutation,
	) (map[string]string, error)
}

type txn struct {
	dgTxn *dgo.Txn
}

func (txn *txn) isCancelErr(err error) bool {
	return status.Code(err) == codes.Canceled
}

func (txn *txn) Query(
	ctx context.Context,
	query string,
	res interface{},
) error {
	rep, err := txn.dgTxn.Query(ctx, query)
	if err != nil {
		if txn.isCancelErr(err) {
			return strerr.New(strerr.ErrCanceled, "")
		}
		return errors.Wrap(err, "query")
	}
	if err := json.Unmarshal(rep.Json, &res); err != nil {
		return errors.Wrap(err, "json unmarsh")
	}
	return nil
}

func (txn *txn) QueryVars(
	ctx context.Context,
	query string,
	vars map[string]string,
	res interface{},
) error {
	rep, err := txn.dgTxn.QueryWithVars(ctx, query, vars)
	if err != nil {
		if txn.isCancelErr(err) {
			return strerr.New(strerr.ErrCanceled, "")
		}
		return errors.Wrap(err, "query")
	}
	if err := json.Unmarshal(rep.Json, &res); err != nil {
		return errors.Wrap(err, "json unmarsh")
	}
	return nil
}

func (txn *txn) Mutation(
	ctx context.Context,
	mutation *api.Mutation,
) (map[string]string, error) {
	assigned, err := txn.dgTxn.Mutate(ctx, mutation)
	if err != nil {
		if txn.isCancelErr(err) {
			return nil, strerr.New(strerr.ErrCanceled, "")
		}
		return nil, errors.Wrap(err, "mutation")
	}
	return assigned.Uids, nil
}

func (str *impl) txn(terr *error) (transaction, func()) {
	// Ensure the database is connected
	if err := str.ensureActive(); err != nil {
		*terr = err
		return nil, nil
	}

	// Create a new transaction and the closure functor
	dgTxn := str.db.NewTxn()
	txn := &txn{
		dgTxn: dgTxn,
	}
	return txn, func() {
		ctx := context.Background()
		if *terr != nil {
			// Rollback transaction
			if rlbErr := dgTxn.Discard(ctx); rlbErr != nil {
				*terr = errors.Wrapf(rlbErr, "rollback after: %s", *terr)
			}
		} else {
			// Commit transaction
			if commitErr := dgTxn.Commit(ctx); commitErr != nil {
				*terr = errors.Wrap(commitErr, "commit")
			}
		}
	}
}
