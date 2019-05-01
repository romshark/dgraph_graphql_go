package apitest

import (
	"flag"
	"fmt"
	"os"
	"testing"

	"github.com/romshark/dgraph_graphql_go/apitest/setup"
)

// stats represents the global statistics recorder the setups must use
var stats = setup.NewStatisticsRecorder()
var dbHost = flag.String("dbhost", "localhost:9080", "database host address")

var tcx setup.TestContext

// TestMain runs the API tests and computes & prints the statistics
func TestMain(m *testing.M) {
	flag.Parse()
	tcx.Stats = stats
	tcx.DBHost = *dbHost

	// Run the tests
	exitCode := m.Run()

	// Compute and print statistics
	computedStats := stats.Compute()
	fmt.Printf("\n Statistics:\n")
	fmt.Printf(
		"  total setups:      %d\n",
		len(computedStats.Tests),
	)
	fmt.Printf(
		"  min setup time:    %s (%s)\n",
		computedStats.MinSetupTime,
		computedStats.MinSetupTimeTest,
	)
	fmt.Printf(
		"  max setup time:    %s (%s)\n",
		computedStats.MaxSetupTime,
		computedStats.MaxSetupTimeTest,
	)
	fmt.Printf(
		"  avg setup time:    %s\n",
		computedStats.AvgSetupTime,
	)
	fmt.Println(" ")
	fmt.Printf(
		"  min teardown time: %s (%s)\n",
		computedStats.MinTeardownTime,
		computedStats.MinTeardownTimeTest,
	)
	fmt.Printf(
		"  max teardown time: %s (%s)\n",
		computedStats.MaxTeardownTime,
		computedStats.MaxTeardownTimeTest,
	)
	fmt.Printf(
		"  avg teardown time: %s\n",
		computedStats.AvgTeardownTime,
	)
	fmt.Println(" ")
	os.Exit(exitCode)
}
