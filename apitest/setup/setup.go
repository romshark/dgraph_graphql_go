package setup

import (
	"context"
	"encoding/base64"
	"testing"
	"time"

	"github.com/romshark/dgraph_graphql_go/api"
)

// TestContext represents a test context
type TestContext struct {
	Stats  *StatisticsRecorder
	DBHost string
}

// TestSetup represents the ArangoDB-based setup of an individual test
type TestSetup struct {
	t                 *testing.T
	stats             *StatisticsRecorder
	apiServer         api.Server
	rootUsername      string
	rootPassword      string
	rootHTTPBasicAuth string
}

// T returns the test reference
func (ts *TestSetup) T() *testing.T {
	return ts.t
}

// New creates a new test setup
func New(t *testing.T, context TestContext) *TestSetup {
	start := time.Now()

	rootUsername := "test"
	rootPassword := "test"

	srvOpts := api.ServerOptions{
		DBHost: context.DBHost,
		RootUser: api.RootUserOptions{
			// Enable the root user in read-write mode
			Status:   api.RootUserRW,
			Username: rootUsername,
			Password: rootPassword,
		},
	}

	apiServer := api.NewServer(srvOpts)
	if err := apiServer.Launch(); err != nil {
		t.Fatalf("API server launch: %s", err)
	}

	testSetup := &TestSetup{
		t:            t,
		stats:        context.Stats,
		apiServer:    apiServer,
		rootUsername: rootUsername,
		rootPassword: rootPassword,
		rootHTTPBasicAuth: "Basic " + base64.StdEncoding.EncodeToString(
			[]byte(rootUsername+":"+rootPassword),
		),
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
