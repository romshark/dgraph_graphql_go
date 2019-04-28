package store

import (
	"context"
	"log"

	"github.com/dgraph-io/dgo"
	"github.com/dgraph-io/dgo/protos/api"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

// Transactions interfaces a transactional store
type Transactions interface {
	CreatePost(
		ctx context.Context,
		author ID,
		title string,
		contents string,
	) error

	CreateReaction(
		ctx context.Context,
		post ID,
		author ID,
		message string,
	) error

	CreateUser(
		ctx context.Context,
		email string,
		displayName string,
	) (ID, error)
}

// Store interfaces a store implementation
type Store interface {
	Prepare() error

	Transactions

	Query(
		ctx context.Context,
		query string,
		result interface{},
	) error

	QueryVars(
		ctx context.Context,
		query string,
		vars map[string]string,
		result interface{},
	) error
}

// store represents the service store
type store struct {
	host    string
	db      *dgo.Dgraph
	onClose func()
}

// NewStore creates a new disconnected database client instance
func NewStore(host string) Store {
	return &store{
		host: host,
		db:   nil,
	}
}

// Prepare prepares the store for use
func (str *store) Prepare() error {
	if str.db != nil {
		return nil
	}

	conn, err := grpc.Dial(str.host, grpc.WithInsecure())
	if err != nil {
		return errors.Wrap(err, "gRPC dial")
	}

	str.db = dgo.NewDgraphClient(api.NewDgraphClient(conn))
	str.onClose = func() {
		if err := conn.Close(); err != nil {
			log.Printf("closing db conn: %s", err)
		}
		str.db = nil
		str.onClose = nil
	}

	return str.setupSchema(context.Background())
}

// IsActive returns true if the store is operational, otherwise returns false
func (str *store) IsActive() bool {
	return str.db != nil
}

func (str *store) ensureActive() error {
	if str.IsActive() {
		return nil
	}
	return errors.New("store inactive")
}
