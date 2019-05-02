package apitest

import (
	"testing"

	"github.com/romshark/dgraph_graphql_go/apitest/setup"
)

// TestCreatePost tests post creation
func TestCreatePost(t *testing.T) {
	ts := setup.New(t, tcx)
	defer ts.Teardown()

	clt := ts.Root()

	author := clt.Help.OK.CreateUser("fooBarowich", "foo@bar.buz", "testpass")
	clt.Help.OK.CreatePost(*author.ID, "test post", "test content")
}
