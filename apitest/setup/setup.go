package setup

import (
	"context"
	"testing"
	"time"

	"github.com/romshark/dgraph_graphql_go/api"
	trn "github.com/romshark/dgraph_graphql_go/api/transport"
	thttp "github.com/romshark/dgraph_graphql_go/api/transport/http"
)

// TestContext represents a test context
type TestContext struct {
	Stats   *StatisticsRecorder
	DBHost  string
	SrvHost string
}

// TestSetup represents the ArangoDB-based setup of an individual test
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

	serverTransport, err := thttp.NewServer(thttp.ServerOptions{
		Host:       context.SrvHost,
		Playground: false,
	})
	if err != nil {
		t.Fatalf("API server transport init: %s", err)
	}

	srvOpts := api.ServerOptions{
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
	if err != nil {
		t.Fatalf("API server init: %s", err)
	}
	if err := apiServer.Launch(); err != nil {
		t.Fatalf("API server launch: %s", err)
	}

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
