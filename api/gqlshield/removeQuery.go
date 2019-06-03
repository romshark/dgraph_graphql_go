package gqlshield

import (
	"fmt"
	"reflect"

	"github.com/pkg/errors"
)

func (shld *shield) RemoveQuery(queryObject Query) error {
	qr, isExpectedType := queryObject.(*query)
	if !isExpectedType {
		return fmt.Errorf(
			"unexpected query type: %s",
			reflect.TypeOf(queryObject),
		)
	}

	shld.lock.Lock()
	defer shld.lock.Unlock()

	if _, deleted := shld.index.Delete(qr.query); !deleted {
		return nil
	}

	deletedQuery := shld.queriesByName[qr.name]
	delete(shld.queriesByName, qr.name)

	if len(qr.query) == shld.longest {
		if err := shld.recalculateLongest(); err != nil {
			return err
		}
	}

	// Persist state changes
	if shld.conf.PersistencyManager != nil {
		if err := shld.conf.PersistencyManager.Save(
			shld.captureState(),
		); err != nil {
			// Rollback changes
			shld.queriesByName[qr.name] = deletedQuery
			shld.index.Insert(qr.query, deletedQuery)
			if err := shld.recalculateLongest(); err != nil {
				rollbackErr := errors.Wrap(
					err,
					"persisting state after removal",
				)
				return errors.Wrap(
					rollbackErr,
					"recalculating longest after rollback",
				)
			}
			return errors.Wrap(err, "persisting state after removal")
		}
	}

	return nil
}
