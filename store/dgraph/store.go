package dgraph

import (
	"context"
	"log"

	"github.com/dgraph-io/dgo"
	"github.com/dgraph-io/dgo/protos/api"
	"github.com/pkg/errors"
	"github.com/romshark/dgraph_graphql_go/api/passhash"
	"github.com/romshark/dgraph_graphql_go/api/sesskeygen"
	"github.com/romshark/dgraph_graphql_go/store"
	"google.golang.org/grpc"
)

// impl represents the service store
type impl struct {
	host                string
	sessionKeyGenerator sesskeygen.SessionKeyGenerator
	passwordHasher      passhash.PasswordHasher
	db                  *dgo.Dgraph
	onClose             func()
}

// NewStore creates a new disconnected database client instance
func NewStore(
	host string,
	sessionKeyGenerator sesskeygen.SessionKeyGenerator,
	passwordHasher passhash.PasswordHasher,
) store.Store {
	return &impl{
		host:                host,
		sessionKeyGenerator: sessionKeyGenerator,
		passwordHasher:      passwordHasher,
		db:                  nil,
	}
}

// Prepare prepares the store for use
func (str *impl) Prepare() error {
	if str.db != nil {
		return nil
	}

	if str.sessionKeyGenerator == nil {
		return errors.Errorf(
			"missing session key generator during store initialization",
		)
	}
	if str.passwordHasher == nil {
		return errors.Errorf(
			"missing password hasher during store initialization",
		)
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
func (str *impl) IsActive() bool {
	return str.db != nil
}

func (str *impl) ensureActive() error {
	if str.IsActive() {
		return nil
	}
	return errors.New("store inactive")
}
