package gqlshield

import (
	"sync"

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

// ClientRole represents a client role
type ClientRole struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// GraphQLShield represents a GraphQL shield instance
type GraphQLShield interface {
	// WhitelistQueries adds the given queries to the whitelist
	// returning an error if one of the queries doesn't meet the requirements.
	WhitelistQueries(newEntry ...Entry) ([]Query, error)

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
		arguments map[string]*string,
	) ([]byte, error)

	// ListQueries returns all whitelisted queries.
	ListQueries() (map[string]Query, error)
}

// NewGraphQLShield creates a new GraphQL shield instance
func NewGraphQLShield(
	config Config,
	clientRoles ...ClientRole,
) (GraphQLShield, error) {
	config.SetDefaults()

	roles, err := copyRoles(clientRoles...)
	if err != nil {
		return nil, err
	}

	shield := &shield{
		conf:          config,
		lock:          &sync.RWMutex{},
		index:         art.New(),
		queriesByName: make(map[string]*query),
		longest:       0,
		clientRoles:   roles,
	}

	if config.PersistencyManager != nil {
		restoredState, err := config.PersistencyManager.Load()
		if err != nil {
			return nil, errors.Wrap(err, "loading state")
		}
		if restoredState == nil {
			// Skip state restoration when no state was loaded
			return shield, nil
		}
		if err := shield.restoreState(restoredState); err != nil {
			return nil, errors.Wrap(err, "restoring state")
		}
	}

	return shield, nil
}

type shield struct {
	conf Config

	// lock synchronizes concurrent access
	lock *sync.RWMutex

	// queriesByName references all whitelisted query objects by their name
	queriesByName map[string]*query

	// index holds a radix-tree lookup index
	index art.Tree

	// longest keeps track of the longest whitelisted query
	longest int

	// clientRoles keeps track of all registered client roles
	clientRoles map[int]ClientRole
}
