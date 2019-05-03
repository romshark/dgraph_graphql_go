package dgraph

import (
	"context"
	"encoding/json"

	"github.com/pkg/errors"
)

func (str *impl) Query(
	ctx context.Context,
	query string,
	result interface{},
) error {
	resp, err := str.db.NewReadOnlyTxn().Query(ctx, query)
	if err != nil {
		return errors.Wrap(err, "query")
	}

	if err := json.Unmarshal(resp.Json, result); err != nil {
		return errors.Wrap(err, "db query result unmarshal")
	}

	return nil
}

func (str *impl) QueryVars(
	ctx context.Context,
	query string,
	vars map[string]string,
	result interface{},
) error {
	resp, err := str.db.NewReadOnlyTxn().QueryWithVars(ctx, query, vars)
	if err != nil {
		return errors.Wrap(err, "query")
	}

	if err := json.Unmarshal(resp.Json, result); err != nil {
		return errors.Wrap(err, "db query result unmarshal")
	}

	return nil
}
