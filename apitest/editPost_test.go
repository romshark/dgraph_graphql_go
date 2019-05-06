package apitest

import (
	"testing"

	"github.com/romshark/dgraph_graphql_go/apitest/setup"
)

// TestEditPost tests post editing
func TestEditPost(t *testing.T) {
	ts := setup.New(t, tcx)
	defer ts.Teardown()

	clt := ts.Root()

	author := clt.Help.OK.CreateUser("fooBarowich", "foo@bar.buz", "testpass")
	post := clt.Help.OK.CreatePost(*author.ID, "test post", "test contents")

	newTitle := "new test post"
	newContents := "new test contents"

	clt.Help.OK.EditPost(
		*post.ID,
		*author.ID,
		&newTitle,
		&newContents,
	)
}
