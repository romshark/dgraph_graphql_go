package apitest

import (
	"testing"

	"github.com/romshark/dgraph_graphql_go/api/graph/gqlmod"
	"github.com/romshark/dgraph_graphql_go/apitest/setup"
	"github.com/romshark/dgraph_graphql_go/store/enum/emotion"
	"github.com/romshark/dgraph_graphql_go/store/errors"
)

// TestCreateReactionAuth tests reaction creation authorization
func TestCreateReactionAuth(t *testing.T) {
	setupTest := func(t *testing.T) (
		ts *setup.TestSetup,
		post *gqlmod.Post,
		commenter *gqlmod.User,
		commenterClt *setup.Client,
	) {
		ts = setup.New(t, tcx)
		debug := ts.Debug()

		author := debug.Help.OK.CreateUser("author", "t@tst.tst", "testpass")
		post = debug.Help.OK.CreatePost(
			*author.ID,
			"test post",
			"test content",
		)
		commenter = debug.Help.OK.CreateUser(
			"commenter",
			"c@tst.tst",
			"testpass",
		)
		commenterClt, _ = ts.Client("c@tst.tst", "testpass")
		return
	}

	// Test creating reactions as a guest
	t.Run("guest (noauth)", func(t *testing.T) {
		ts, post, cmt, _ := setupTest(t)
		defer ts.Teardown()

		ts.Guest().Help.ERR.CreateReaction(
			errors.ErrUnauthorized,
			*cmt.ID,
			*post.ID,
			emotion.Excited,
			"test comment",
		)
	})

	// Test creating reactions on behalf of other users
	t.Run("non-author (noauth)", func(t *testing.T) {
		ts, post, _, cmtClt := setupTest(t)
		defer ts.Teardown()

		other := ts.Debug().Help.OK.CreateUser("other", "t2@tst.tst", "testpass")
		cmtClt.Help.ERR.CreateReaction(
			errors.ErrUnauthorized,
			*other.ID, // Different reaction author ID
			*post.ID,
			emotion.Excited,
			"test comment",
		)
	})
}
