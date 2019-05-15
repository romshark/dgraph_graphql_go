package apitest

import (
	"testing"

	"github.com/romshark/dgraph_graphql_go/apitest/setup"
)

// TestCreatePost tests post creation
func TestCreatePost(t *testing.T) {
	ts := setup.New(t, tcx)
	defer ts.Teardown()

	root := ts.Root()

	author := root.Help.OK.CreateUser("fooBarowich", "foo@bar.buz", "testpass")
	root.Help.OK.CreatePost(*author.ID, "test post", "test content")
}
