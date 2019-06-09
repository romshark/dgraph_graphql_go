package gqlshield

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
)

func (shld *shield) WhitelistQueries(newEntries ...Entry) ([]Query, error) {
	// Create new normalized and prepared query instances
	newQueries := make([]*query, len(newEntries))
	for i, newEntry := range newEntries {
		newQuery := &query{
			id:       newID(),
			creation: time.Now(),
		}

		// Ensure query name validity
		if err := validateQueryName(newEntry.Name); err != nil {
			return nil, err
		}

		// Set name
		newQuery.name = newEntry.Name

		// Ensure query string validity
		if err := validateQueryString(newEntry.Query); err != nil {
			return nil, err
		}

		// Set query (normalized)
		newQuery.query = []byte(newEntry.Query)
		normalized, err := prepareQuery(newQuery.query)
		if err != nil {
			return nil, err
		}
		newQuery.query = normalized

		// Ensure whitelistedFor validity
		if len(newEntry.WhitelistedFor) < 1 {
			return nil, fmt.Errorf(
				"query '%s' has no roles associated",
				newEntry.Name,
			)
		}

		// Set whitelistedFor
		newQuery.whitelistedFor = make(
			map[int]struct{},
			len(newEntry.WhitelistedFor),
		)
		for _, roleID := range newEntry.WhitelistedFor {
			// Ensure whitelistedFor role ID uniqueness
			if _, isDefined := newQuery.whitelistedFor[roleID]; isDefined {
				return nil, fmt.Errorf(
					"query '%s' has duplicate role IDs (%d) in whitelistedFor",
					newEntry.Name,
					roleID,
				)
			}

			newQuery.whitelistedFor[roleID] = struct{}{}
		}

		// Set parameters
		if newEntry.Parameters != nil {
			newQuery.parameters = make(
				map[string]Parameter,
				len(newEntry.Parameters),
			)
			for paramName, param := range newEntry.Parameters {
				// Ensure parameter name validity
				if err := validateParameterName(paramName); err != nil {
					return nil, errors.Wrapf(
						err,
						"query '%s' has parameter with invalid name",
						newEntry.Name,
					)
				}

				// Ensure parameter properties validity
				if param.MaxValueLength < 1 {
					return nil, fmt.Errorf(
						"query '%s' has parameter ('%s') with invalid "+
							"property MaxValueLength (%d)",
						newEntry.Name,
						paramName,
						param.MaxValueLength,
					)
				}

				// Ensure parameter name uniqueness
				if _, isDefined := newQuery.parameters[paramName]; isDefined {
					return nil, fmt.Errorf(
						"query '%s' has duplicate parameter name ('%s')",
						newEntry.Name,
						paramName,
					)
				}

				newQuery.parameters[paramName] = param
			}
		}

		newQueries[i] = newQuery
	}

	queries := make([]Query, len(newQueries))
	for i, newQuery := range newQueries {
		queries[i] = newQuery
	}

	// Verify queries against the store
	shld.lock.Lock()
	defer shld.lock.Unlock()

	for i, newQuery := range newQueries {
		// Ensure referenced roles exist
		for role := range newQuery.whitelistedFor {
			if _, roleDefined := shld.clientRoles[role]; !roleDefined {
				return nil, fmt.Errorf("undefined role: %d", role)
			}
		}

		// Ensure name uniqueness
		if _, nameExists := shld.queriesByName[newQuery.name]; nameExists {
			return nil, errors.Errorf(
				"%d: a query with a similar name (%s) is already whitelisted",
				i,
				newQuery.name,
			)
		}

		// Ensure query uniqueness
		if existing, similarQueryRegistered := shld.index.Search(
			newQuery.query,
		); similarQueryRegistered {
			return nil, fmt.Errorf(
				"similar query already whitelisted under the name: '%s'",
				existing.(*query).name,
			)
		}
	}

	// Update state
	for _, newQuery := range newQueries {
		// Store the original query
		shld.queriesByName[newQuery.name] = newQuery

		// Update index
		shld.index.Insert(newQuery.query, newQuery)
		if len(newQuery.query) > shld.longest {
			shld.longest = len(newQuery.query)
		}

		// Persist state changes
		if shld.conf.PersistencyManager != nil {
			if err := shld.conf.PersistencyManager.Save(
				shld.captureState(),
			); err != nil {
				// Rollback changes
				delete(shld.queriesByName, newQuery.name)
				shld.index.Delete(newQuery.query)
				if err := shld.recalculateLongest(); err != nil {
					rollbackErr := errors.Wrap(
						err,
						"persisting state after insertion",
					)
					return nil, errors.Wrap(
						rollbackErr,
						"recalculating longest after rollback",
					)
				}

				return nil, errors.Wrap(err, "persisting state after insertion")
			}
		}
	}

	return queries, nil
}
