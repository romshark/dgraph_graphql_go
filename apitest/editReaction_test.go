package apitest

import (
	"testing"

	"github.com/romshark/dgraph_graphql_go/apitest/setup"
	"github.com/romshark/dgraph_graphql_go/store/enum/emotion"
)

// TestEditReaction tests reaction editing
func TestEditReaction(t *testing.T) {
	ts := setup.New(t, tcx)
	defer ts.Teardown()

	// Prepare
	debug := ts.Debug()
	author := debug.Help.OK.CreateUser("author", "1@tst.tst", "testpass")
	authorClt, _ := ts.Client("1@tst.tst", "testpass")
	post := debug.Help.OK.CreatePost(*author.ID, "test post", "test contents")
	reaction := debug.Help.OK.CreateReaction(
		*author.ID,
		*post.ID,
		emotion.Happy,
		"sample message",
	)

	// Test edit
	authorClt.Help.OK.EditReaction(
		*reaction.ID,
		*author.ID,
		"new message",
	)
}
