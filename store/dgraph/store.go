package dgraph

import (
	"context"
	"log"

	"github.com/dgraph-io/dgo/v2"
	"github.com/dgraph-io/dgo/v2/protos/api"
	"github.com/pkg/errors"
	"github.com/romshark/dgraph_graphql_go/store"
	"google.golang.org/grpc"
)

// impl represents the service store
type impl struct {
	host            string
	db              *dgo.Dgraph
	comparePassword func(hash, password string) bool
	onClose         func()
	debugLog        *log.Logger
	errorLog        *log.Logger
}

// NewStore creates a new disconnected database client instance
func NewStore(
	host string,
	comparePassword func(hash, password string) bool,
	debugLog *log.Logger,
	errorLog *log.Logger,
) store.Store {
	return &impl{
		host:            host,
		db:              nil,
		comparePassword: comparePassword,
		debugLog:        debugLog,
		errorLog:        errorLog,
	}
}

// Prepare prepares the store for use
func (str *impl) Prepare() error {
	if str.db != nil {
		return nil
	}

	conn, err := grpc.Dial(str.host, grpc.WithInsecure())
	if err != nil {
		return errors.Wrap(err, "gRPC dial")
	}
	str.debugLog.Printf("database (%s) connected", str.host)

	str.db = dgo.NewDgraphClient(api.NewDgraphClient(conn))
	str.onClose = func() {
		if err := conn.Close(); err != nil {
			str.errorLog.Printf("closing db conn: %s", err)
		}
		str.db = nil
		str.onClose = nil
	}

	err = str.setupSchema(context.Background())
	if err != nil {
		str.debugLog.Print("database schema setup")
	}

	return err
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
