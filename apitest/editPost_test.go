package apitest

import (
	"testing"

	"github.com/romshark/dgraph_graphql_go/apitest/setup"
)

// TestEditPost tests post editing
func TestEditPost(t *testing.T) {
	ts := setup.New(t, tcx)
	defer ts.Teardown()

	// Prepare
	debug := ts.Debug()
	author := debug.Help.OK.CreateUser("author", "1@tst.tst", "testpass")
	authorClt, _ := ts.Client("1@tst.tst", "testpass")
	post := debug.Help.OK.CreatePost(*author.ID, "test post", "test contents")

	// Test edit
	newTitle := "new test post"
	newContents := "new test contents"
	authorClt.Help.OK.EditPost(
		*post.ID,
		*author.ID,
		&newTitle,
		&newContents,
	)
}
