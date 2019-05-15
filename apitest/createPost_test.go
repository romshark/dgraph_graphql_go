package apitest

import (
	"testing"

	"github.com/romshark/dgraph_graphql_go/apitest/setup"
)

// TestCreatePost tests post creation
func TestCreatePost(t *testing.T) {
	ts := setup.New(t, tcx)
	defer ts.Teardown()

	author := ts.Debug().Help.OK.CreateUser("usr1", "t@tst.tst", "testpass")
	authorClt, _ := ts.Client("t@tst.tst", "testpass")

	authorClt.Help.OK.CreatePost(*author.ID, "test post", "test content")
}
