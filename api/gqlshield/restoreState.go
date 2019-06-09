package gqlshield

import (
	"github.com/pkg/errors"
	art "github.com/plar/go-adaptive-radix-tree"
)

func (shld *shield) restoreState(state *State) error {
	// Restore roles
	clientRoles := make(map[int]ClientRole, len(state.Roles))
	findRoleByName := func(name string) *int {
		for id, role := range clientRoles {
			if name == role.Name {
				return &id
			}
		}
		return nil
	}
	for _, role := range state.Roles {
		// Ensure role ID uniqueness
		if _, alreadyDefined := clientRoles[role.ID]; alreadyDefined {
			return errors.Errorf("duplicate role IDs (%d)", role.ID)
		}
		// Ensure role name uniqueness
		if foundID := findRoleByName(role.Name); foundID != nil {
			return errors.Errorf(
				"duplicate role name (%d and %d)",
				role.ID,
				foundID,
			)
		}
		clientRoles[role.ID] = role
	}

	// Restore queries
	index := art.New()
	queriesByName := make(map[string]*query)

	for id, queryModel := range state.WhitelistedQueries {
		// Ensure the query ID is valid
		identifier := ID(id)
		if err := identifier.Validate(); err != nil {
			return err
		}

		// Verify query name validity
		if err := validateQueryName(queryModel.Name); err != nil {
			return errors.Wrapf(err, "query %s has invalid name", id)
		}

		// Ensure query name uniqueness
		if _, used := queriesByName[queryModel.Name]; used {
			return errors.Errorf(
				"duplicate query name ('%s')",
				queryModel.Name,
			)
		}

		// Normalize query string
		queryString := []byte(queryModel.Query)
		var err error
		queryString, err = prepareQuery(queryString)
		if err != nil {
			return errors.Wrap(err, "preparing query")
		}

		// Verify query string validity
		if err := validateQueryString(queryModel.Name); err != nil {
			return errors.Wrapf(err, "query %s has invalid query string", id)
		}

		// Ensure query string uniqueness
		if _, found := index.Search(queryString); found {
			return errors.Errorf("duplicate query ('%s')", queryModel.Query)
		}

		// Verify parameters
		for paramName, param := range queryModel.Parameters {
			if err := validateParameterName(paramName); err != nil {
				return errors.Wrapf(
					err,
					"query %s has invalid parameter name",
					id,
				)
			}
			if param.MaxValueLength < 1 {
				return errors.Errorf(
					"query %s has parameter ('%s') with invalid "+
						"property MaxValueLength (%d)",
					id,
					paramName,
					param.MaxValueLength,
				)
			}
		}

		// Verify the list of role IDs the query is whitelisted for
		if len(queryModel.WhitelistedFor) < 1 {
			return errors.Errorf("query %s has no role IDs associated", id)
		}
		whitelistedFor := make(map[int]struct{}, len(queryModel.WhitelistedFor))
		for _, roleID := range queryModel.WhitelistedFor {
			// Ensure the referenced role ID is defined
			if _, isDefined := clientRoles[roleID]; !isDefined {
				return errors.Errorf("undefined role ID (%d)", roleID)
			}

			// Ensure there are no duplicate role IDs in WhitelistedFor
			if _, alreadyDefined := whitelistedFor[roleID]; alreadyDefined {
				return errors.Errorf(
					"query %s contains duplicate role IDs (%d)",
					id,
					roleID,
				)
			}

			whitelistedFor[roleID] = struct{}{}
		}

		query := &query{
			id:             ID(id),
			query:          []byte(queryModel.Query),
			creation:       queryModel.Creation,
			name:           queryModel.Name,
			parameters:     queryModel.Parameters,
			whitelistedFor: whitelistedFor,
		}

		queriesByName[queryModel.Name] = query
		index.Insert(queryString, query)
	}

	shld.lock.Lock()
	defer shld.lock.Unlock()

	oldQueriesByName := shld.queriesByName
	oldIndex := shld.index

	shld.queriesByName = queriesByName
	shld.index = index

	if err := shld.recalculateLongest(); err != nil {
		// Rollback changes
		shld.queriesByName = oldQueriesByName
		shld.index = oldIndex
		return errors.Wrap(err, "recalculating longest")
	}

	return nil
}
