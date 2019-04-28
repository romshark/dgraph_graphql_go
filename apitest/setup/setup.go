package setup

import (
	"context"
	"demo/api"
	"demo/store"
	"net/http"
	"testing"
	"time"
)

// TestContext represents a test context
type TestContext struct {
	Stats  *StatisticsRecorder
	DBHost string
}

// TestSetup represents the ArangoDB-based setup of an individual test
type TestSetup struct {
	t             *testing.T
	stats         *StatisticsRecorder
	apiServer     api.Server
	defaultClient *http.Client
}

// New creates a new test setup
func New(t *testing.T, context TestContext) *TestSetup {
	start := time.Now()

	// Launch API server
	str := store.NewStore(context.DBHost)

	apiServer := api.NewServer(api.ServerOptions{}, str)
	if err := apiServer.Launch(); err != nil {
		t.Fatalf("API server launch: %s", err)
	}

	testSetup := &TestSetup{
		t:         t,
		stats:     context.Stats,
		apiServer: apiServer,
		defaultClient: &http.Client{
			Timeout: time.Second * 1,
		},
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
