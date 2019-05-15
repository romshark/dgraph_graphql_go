package apitest

import (
	"testing"

	"github.com/romshark/dgraph_graphql_go/apitest/setup"
)

// TestSignIn tests sign in
func TestSignIn(t *testing.T) {
	ts := setup.New(t, tcx)
	defer ts.Teardown()

	root := ts.Root()
	root.Help.OK.CreateUser("fooBarowich", "foo@bar.buz", "testpass")

	guest := ts.Guest()
	guest.Help.OK.SignIn("foo@bar.buz", "testpass")
}
