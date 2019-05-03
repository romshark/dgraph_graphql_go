package store

import (
	"context"
	"log"
	"time"

	"github.com/romshark/dgraph_graphql_go/store/enum/emotion"

	"github.com/dgraph-io/dgo"
	"github.com/dgraph-io/dgo/protos/api"
	"github.com/pkg/errors"
	"github.com/romshark/dgraph_graphql_go/api/passhash"
	"github.com/romshark/dgraph_graphql_go/api/sesskeygen"
	"google.golang.org/grpc"
)

// Transactions interfaces a transactional store
type Transactions interface {
	CreateSession(
		ctx context.Context,
		email string,
		password string,
	) (
		uid UID,
		key string,
		creation time.Time,
		userUid UID,
		err error,
	)

	CreatePost(
		ctx context.Context,
		author ID,
		title string,
		contents string,
	) (UID, ID, error)

	CreateReaction(
		ctx context.Context,
		author ID,
		subject ID,
		emotion emotion.Emotion,
		message string,
	) (
		result struct {
			UID          UID
			ID           ID
			SubjectUID   UID
			AuthorUID    UID
			CreationTime time.Time
		},
		err error,
	)

	CreateUser(
		ctx context.Context,
		email string,
		displayName string,
		password string,
	) (UID, ID, error)
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
) Store {
	return &store{
		host:                host,
		sessionKeyGenerator: sessionKeyGenerator,
		passwordHasher:      passwordHasher,
		db:                  nil,
	}
}

// Prepare prepares the store for use
func (str *store) Prepare() error {
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
func (str *store) IsActive() bool {
	return str.db != nil
}

func (str *store) ensureActive() error {
	if str.IsActive() {
		return nil
	}
	return errors.New("store inactive")
}
