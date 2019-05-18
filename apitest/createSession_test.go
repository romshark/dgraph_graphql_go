package apitest

import (
	"testing"

	"github.com/romshark/dgraph_graphql_go/apitest/setup"
)

// TestCreateSession tests session creation
func TestCreateSession(t *testing.T) {
	ts := setup.New(t, tcx)
	defer ts.Teardown()

	ts.Debug().Help.OK.CreateUser("fooBarowich", "foo@bar.buz", "testpass")

	guest := ts.Guest()
	guest.Help.OK.CreateSession("foo@bar.buz", "testpass")
}
