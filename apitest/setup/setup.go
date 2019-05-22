package setup

import (
	"context"
	ctx "context"
	"testing"
	"time"

	"github.com/dgraph-io/dgo"
	dbapi "github.com/dgraph-io/dgo/protos/api"
	"github.com/romshark/dgraph_graphql_go/api"
	trn "github.com/romshark/dgraph_graphql_go/api/transport"
	thttp "github.com/romshark/dgraph_graphql_go/api/transport/http"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

// TestContext represents a test context
type TestContext struct {
	Stats   *StatisticsRecorder
	DBHost  string
	SrvHost string
}

// TestSetup represents the Dgraph-based server setup of an individual test
type TestSetup struct {
	t               *testing.T
	stats           *StatisticsRecorder
	apiServer       api.Server
	serverTransport trn.Server
	debugUsername   string
	debugPassword   string
}

// T returns the test reference
func (ts *TestSetup) T() *testing.T {
	return ts.t
}

// New creates a new test setup
func New(t *testing.T, context TestContext) *TestSetup {
	start := time.Now()

	debugUsername := "test"
	debugPassword := "test"

	// Clear database
	conn, err := grpc.Dial(context.DBHost, grpc.WithInsecure())
	require.NoError(t, err)
	db := dgo.NewDgraphClient(dbapi.NewDgraphClient(conn))
	require.NoError(t, db.Alter(
		ctx.Background(),
		&dbapi.Operation{DropAll: true},
	))
	require.NoError(t, conn.Close())

	serverTransport, err := thttp.NewServer(thttp.ServerOptions{
		Host:       context.SrvHost,
		Playground: false,
	})
	require.NoError(t, err)

	srvOpts := api.ServerOptions{
		Mode:   api.ModeDebug,
		DBHost: context.DBHost,
		DebugUser: api.DebugUserOptions{
			// Enable the debug user in read-write mode
			Status:   api.DebugUserRW,
			Username: debugUsername,
			Password: debugPassword,
		},
		Transport: []trn.Server{
			serverTransport,
		},
	}

	apiServer, err := api.NewServer(srvOpts)
	require.NoError(t, err)
	require.NoError(t, apiServer.Launch())

	testSetup := &TestSetup{
		t:               t,
		stats:           context.Stats,
		apiServer:       apiServer,
		serverTransport: serverTransport,
		debugUsername:   debugUsername,
		debugPassword:   debugPassword,
	}

	// Record setup time
	context.Stats.Set(t, func(stat *TestStatistics) {
		stat.SetupTime = time.Since(start)
	})

	return testSetup
}

// Teardown gracefully terminates the test,
// this method MUST BE DEFERRED until the end of the test!
func (ts *TestSetup) Teardown() {
	start := time.Now()

	// Stop the API server instance
	if err := ts.apiServer.Shutdown(context.Background()); err != nil {
		// Don't break on shutdown failure, remove database before quitting!
		ts.t.Errorf("API server shutdown: %s", err)
	}

	// Record teardown time
	ts.stats.Set(ts.t, func(stat *TestStatistics) {
		stat.TeardownTime = time.Since(start)
	})
}
