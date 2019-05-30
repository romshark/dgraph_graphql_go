package gqlshield

import (
	"errors"
	"sync"

	art "github.com/plar/go-adaptive-radix-tree"
)

// Parameter represents a query parameter
type Parameter struct {
	Name           string
	MaxValueLength uint32
}

// Query represents a whitelisted query
type Query struct {
	Query      []byte
	Name       string
	Parameters map[string]Parameter
}

// Clone creates a new exact copy of the query object
func (q *Query) Clone() *Query {
	qr := make([]byte, len(q.Query))
	copy(qr, q.Query)

	params := make(map[string]Parameter, len(q.Parameters))
	for paramName, param := range q.Parameters {
		params[paramName] = param
	}

	return &Query{
		Query:      qr,
		Name:       q.Name,
		Parameters: params,
	}
}

// GraphQLShield represents a GraphQL shield instance
type GraphQLShield interface {
	// WhitelistQuery adds a query to the whitelist
	WhitelistQuery(query *Query) error

	// RemoveQuery removes a query from the whitelist and returns true
	// if any query was removed
	RemoveQuery(query []byte) (*Query, error)

	// Check returns an error if the given query isn't whitelisted
	Check(query []byte) (bool, error)

	// Queries returns all whitelisted queries
	Queries() map[string]*Query
}

// NewGraphQLShield creates a new GraphQL shield instance
func NewGraphQLShield() GraphQLShield {
	return &shield{
		lock:  &sync.RWMutex{},
		index: art.New(),
		store: make(map[string]*Query),
	}
}

type shield struct {
	lock *sync.RWMutex

	// store stores all original queries associted by a name
	store map[string]*Query

	// index holds a radix-tree lookup index
	index art.Tree
}

func (shld *shield) WhitelistQuery(query *Query) error {
	if len(query.Query) < 1 {
		return errors.New("invalid (empty) query")
	}
	normalized, err := prepareQuery(query.Query)
	if err != nil {
		return err
	}

	shld.lock.Lock()
	defer shld.lock.Unlock()
	if _, nameIsRegistered := shld.store[query.Name]; nameIsRegistered {
		return errors.New("a query with a similar name is already whitelisted")
	}

	// Store the original query
	shld.store[query.Name] = query

	shld.index.Insert(normalized, query)
	return nil
}

func (shld *shield) RemoveQuery(query []byte) (*Query, error) {
	if len(query) < 1 {
		return nil, errors.New("invalid (empty) query")
	}
	normalized, err := prepareQuery(query)
	if err != nil {
		return nil, err
	}

	shld.lock.Lock()
	defer shld.lock.Unlock()

	val, deleted := shld.index.Delete(normalized)
	if deleted {
		removed := val.(*Query)
		delete(shld.store, removed.Name)
		return removed, nil
	}

	return nil, nil
}

func (shld *shield) Check(query []byte) (bool, error) {
	if len(query) < 1 {
		return false, errors.New("invalid (empty) query")
	}
	normalized, err := prepareQuery(query)
	if err != nil {
		return false, err
	}

	shld.lock.RLock()
	defer shld.lock.RUnlock()

	_, found := shld.index.Search(normalized)
	return found, nil
}

func (shld *shield) Queries() map[string]*Query {
	shld.lock.RLock()
	defer shld.lock.RUnlock()

	m := make(map[string]*Query, shld.index.Size())
	for itr := shld.index.Iterator(); itr.HasNext(); {
		node, _ := itr.Next()
		qr := node.Value().(*Query)
		m[qr.Name] = qr
	}
	return m
}
