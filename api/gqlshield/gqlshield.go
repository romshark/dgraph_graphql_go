package gqlshield

import (
	"fmt"
	"reflect"
	"sync"
	"time"

	"github.com/pkg/errors"
	art "github.com/plar/go-adaptive-radix-tree"
)

// Entry represents a whitelist entry prototype
type Entry struct {
	Query          string
	Name           string
	Parameters     map[string]Parameter
	WhitelistedFor []int
}

// GraphQLShield represents a GraphQL shield instance
type GraphQLShield interface {
	// WhitelistQuery adds the given query to the whitelist
	// returning an error if the query doesn't meet the requirements.
	WhitelistQuery(newEntry Entry) (Query, error)

	// RemoveQuery removes a query from the whitelist and returns true
	// if any query was removed as well as the actual removed query.
	RemoveQuery(query Query) error

	// Check returns an error if the given query isn't allowed for the given
	// client role to be executed or if the provided arguments are unacceptable.
	//
	// WARNING: query will be mutated during normalization! Manually copy the
	// query byte-slice if you don't want your inputs to be mutated.
	Check(
		clientRole int,
		query []byte,
		arguments map[string]string,
	) error

	// Queries returns all whitelisted queries.
	Queries() (map[string]Query, error)
}

// ClientRole represents a client role
type ClientRole struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func copyRoles(roles ...ClientRole) (map[int]ClientRole, error) {
	if len(roles) < 1 {
		return nil, errors.New("missing roles")
	}

	byName := make(map[string]ClientRole, len(roles))
	byID := make(map[int]ClientRole, len(roles))
	cp := make(map[int]ClientRole, len(roles))
	for _, role := range roles {
		// Ensure the client role ID is unique
		if defined, alreadyDefined := byID[role.ID]; alreadyDefined {
			return nil, fmt.Errorf(
				"role ID %d redefined for %s",
				defined.ID,
				defined.Name,
			)
		}

		// Ensure the client role name is valid
		if len(role.Name) < 1 {
			return nil, fmt.Errorf(
				"invalid (empty) client role name (%d)",
				role.ID,
			)
		}

		// Ensure the client role name is unique
		if defined, alreadyDefined := byName[role.Name]; alreadyDefined {
			return nil, fmt.Errorf(
				"%d:'%s' redefined for %d",
				defined.ID,
				role.Name,
				role.ID,
			)
		}

		role := ClientRole{
			ID:   role.ID,
			Name: role.Name,
		}

		// Define name->id
		byName[role.Name] = role
		byID[role.ID] = role
		cp[role.ID] = role
	}
	return cp, nil
}

// NewGraphQLShield creates a new GraphQL shield instance
func NewGraphQLShield(
	config Config,
	clientRoles ...ClientRole,
) (GraphQLShield, error) {
	roles, err := copyRoles(clientRoles...)
	if err != nil {
		return nil, err
	}

	if config.PersistencyManager != nil {
		config.PersistencyManager.Load()
	}

	return &shield{
		conf:        config,
		lock:        &sync.RWMutex{},
		index:       art.New(),
		store:       make(map[string]*query),
		longest:     0,
		clientRoles: roles,
	}, nil
}

type shield struct {
	conf Config

	// lock synchronizes concurrent access
	lock *sync.RWMutex

	// store stores all original queries associted by a name
	store map[string]*query

	// index holds a radix-tree lookup index
	index art.Tree

	// longest keeps track of the longest whitelisted query
	longest int

	// clientRoles keeps track of all registered client roles
	clientRoles map[int]ClientRole
}

func (shld *shield) captureState() *State {
	roles := make([]ClientRole, 0, len(shld.clientRoles))
	for _, role := range shld.clientRoles {
		roles = append(roles, role)
	}

	queries := make(map[string]QueryModel, len(shld.store))
	for _, query := range shld.store {
		var params map[string]Parameter
		if query.parameters != nil {
			params = make(map[string]Parameter, len(query.parameters))
			for name, param := range query.parameters {
				params[name] = param
			}
		}

		whitelistedFor := make([]int, 0, len(query.whitelistedFor))
		for role := range query.whitelistedFor {
			whitelistedFor = append(whitelistedFor, role)
		}

		queries[string(query.id)] = QueryModel{
			Query:          string(query.query),
			Creation:       query.creation,
			Name:           query.name,
			Parameters:     params,
			WhitelistedFor: whitelistedFor,
		}
	}

	return &State{
		Roles:              roles,
		WhitelistedQueries: queries,
	}
}

func (shld *shield) recalculateLongest() error {
	// Recalculate longest
	shld.longest = 0
	for itr := shld.index.Iterator(); itr.HasNext(); {
		node, err := itr.Next()
		if err != nil {
			return err
		}
		queryLength := len(node.Value().(*query).query)
		if queryLength > shld.longest {
			shld.longest = queryLength
		}
	}
	return nil
}

func (shld *shield) WhitelistQuery(newEntry Entry) (Query, error) {
	newQuery := &query{
		id:       newID(),
		creation: time.Now(),
	}

	// Set name
	if len(newEntry.Name) < 1 {
		return nil, errors.New("invalid (empty) query name")
	}
	newQuery.name = newEntry.Name

	// Set query (normalized)
	if len(newEntry.Query) < 1 {
		return nil, errors.New("invalid (empty) query")
	}
	newQuery.query = []byte(newEntry.Query)
	normalized, err := prepareQuery(newQuery.query)
	if err != nil {
		return nil, err
	}
	newQuery.query = normalized

	// Set whitelistedFor
	if len(newEntry.WhitelistedFor) < 1 {
		return nil, fmt.Errorf(
			"query '%s' has no roles associated",
			newEntry.Name,
		)
	}
	newQuery.whitelistedFor = make(
		map[int]struct{},
		len(newEntry.WhitelistedFor),
	)
	for _, roleID := range newEntry.WhitelistedFor {
		newQuery.whitelistedFor[roleID] = struct{}{}
	}

	// Set parameters
	if newEntry.Parameters != nil {
		newQuery.parameters = make(
			map[string]Parameter,
			len(newEntry.Parameters),
		)
		for paramName, param := range newEntry.Parameters {
			if len(paramName) < 1 {
				return nil, fmt.Errorf(
					"query '%s' has parameter with invalid (empty name)",
					newEntry.Name,
				)
			}
			if param.MaxValueLength < 1 {
				return nil, fmt.Errorf(
					"query '%s' has parameter with invalid MaxValueLength: %d",
					newEntry.Name,
					param.MaxValueLength,
				)
			}
			newQuery.parameters[paramName] = param
		}
	}

	// Verify & register query
	shld.lock.Lock()
	defer shld.lock.Unlock()

	// Ensure name uniqueness
	if _, nameIsRegistered := shld.store[newEntry.Name]; nameIsRegistered {
		return nil, errors.New(
			"a query with a similar name is already whitelisted",
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

	// Ensure referenced roles exist
	for role := range newQuery.whitelistedFor {
		if _, roleDefined := shld.clientRoles[role]; !roleDefined {
			return nil, fmt.Errorf("undefined role: %d", role)
		}
	}

	/* Update state */

	// Store the original query
	shld.store[newEntry.Name] = newQuery

	// Update index
	shld.index.Insert(normalized, newQuery)
	if len(newQuery.query) > shld.longest {
		shld.longest = len(newQuery.query)
	}

	// Persist state changes
	if shld.conf.PersistencyManager != nil {
		if err := shld.conf.PersistencyManager.Save(
			shld.captureState(),
		); err != nil {
			// Rollback changes
			delete(shld.store, newEntry.Name)
			shld.index.Delete(normalized)
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

	return newQuery, nil
}

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

	deletedQuery := shld.store[qr.name]
	delete(shld.store, qr.name)

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
			shld.store[qr.name] = deletedQuery
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

func (shld *shield) Check(
	clientRoleID int,
	Query []byte,
	arguments map[string]string,
) error {
	if len(Query) < 1 {
		return Error{
			Code:    ErrWrongInput,
			Message: "invalid (empty) query",
		}
	}
	normalized, err := prepareQuery(Query)
	if err != nil {
		return err
	}

	shld.lock.RLock()
	defer shld.lock.RUnlock()

	// Find role
	if _, roleDefined := shld.clientRoles[clientRoleID]; !roleDefined {
		return fmt.Errorf("role %d is undefined", clientRoleID)
	}

	// Lookup query
	qrObj, found := shld.index.Search(normalized)
	if !found {
		return Error{
			Code:    ErrUnauthorized,
			Message: "query not whitelisted",
		}
	}
	qr := qrObj.(*query)

	// Ensure the client is allowed to execute this query
	if _, roleAllowed := qr.whitelistedFor[clientRoleID]; !roleAllowed {
		return Error{
			Code: ErrUnauthorized,
			Message: fmt.Sprintf(
				"role %d is not allowed to execute this query",
				clientRoleID,
			),
		}
	}

	// Check arguments
	if len(arguments) != len(qr.parameters) {
		return Error{
			Code: ErrUnauthorized,
			Message: fmt.Sprintf(
				"unexpected number of arguments: (%d/%d)",
				len(arguments),
				len(qr.parameters),
			),
		}
	}
	for name, expectedParam := range qr.parameters {
		actual, hasArg := arguments[name]
		if !hasArg {
			return Error{
				Code:    ErrUnauthorized,
				Message: fmt.Sprintf("missing argument '%s'", name),
			}
		}
		if uint32(len(actual)) > expectedParam.MaxValueLength {
			return Error{
				Code: ErrUnauthorized,
				Message: fmt.Sprintf(
					"argument '%s' exceeds max length (%d/%d)",
					name,
					len(actual),
					expectedParam.MaxValueLength,
				),
			}
		}
	}

	return nil
}

func (shld *shield) Queries() (map[string]Query, error) {
	shld.lock.RLock()
	defer shld.lock.RUnlock()

	allQueries := make(map[string]Query, shld.index.Size())
	for itr := shld.index.Iterator(); itr.HasNext(); {
		node, err := itr.Next()
		if err != nil {
			return nil, err
		}
		qr := node.Value().(*query)
		allQueries[qr.name] = qr
	}
	return allQueries, nil
}
