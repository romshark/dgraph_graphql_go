package apitest

import (
	"testing"

	"github.com/romshark/dgraph_graphql_go/apitest/setup"
)

// TestCreateUser tests user account creation
func TestCreateUser(t *testing.T) {
	ts := setup.New(t, tcx)
	defer ts.Teardown()

	ts.Help.OK.CreateUser("fooBarowich", "foo@bar.buz")
}
