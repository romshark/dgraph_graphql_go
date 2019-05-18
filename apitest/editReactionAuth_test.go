package apitest

import (
	"testing"

	"github.com/romshark/dgraph_graphql_go/api/graph/gqlmod"
	"github.com/romshark/dgraph_graphql_go/apitest/setup"
	"github.com/romshark/dgraph_graphql_go/store/enum/emotion"
	"github.com/romshark/dgraph_graphql_go/store/errors"
)

// TestEditReactionAuth tests reaction editing authorization
func TestEditReactionAuth(t *testing.T) {
	setupTest := func(t *testing.T) (
		ts *setup.TestSetup,
		author *gqlmod.User,
		authorClt *setup.Client,
		reaction *gqlmod.Reaction,
	) {
		ts = setup.New(t, tcx)
		debug := ts.Debug()

		authorEmail := "author@tst.tst"
		authorPass := "testpass"
		author = debug.Help.OK.CreateUser(
			"fooBarowich",
			authorEmail,
			authorPass,
		)
		authorClt, _ = ts.Client(authorEmail, authorPass)

		post := debug.Help.OK.CreatePost(
			*author.ID,
			"example title",
			"example contents",
		)
		reaction = debug.Help.OK.CreateReaction(
			*author.ID,
			*post.ID,
			emotion.Happy,
			"sample message",
		)

		return
	}

	// Test editing reactions as a guest
	t.Run("guest (noauth)", func(t *testing.T) {
		ts, author, _, reaction := setupTest(t)
		defer ts.Teardown()

		ts.Guest().Help.ERR.EditReaction(
			errors.ErrUnauthorized,
			*reaction.ID,
			*author.ID,
			"new message",
		)
	})

	// Test editing reactions on behalf of other users
	t.Run("non-editor (noauth)", func(t *testing.T) {
		ts, _, authorClt, reaction := setupTest(t)
		defer ts.Teardown()

		other := ts.Debug().Help.OK.CreateUser("other", "2@tst.tst", "testpass")

		authorClt.Help.ERR.EditReaction(
			errors.ErrUnauthorized,
			*reaction.ID,
			*other.ID, // Different editor ID
			"new message",
		)
	})

	// Test editing reactions of other users
	t.Run("non-author (noauth)", func(t *testing.T) {
		ts, _, _, reaction := setupTest(t)
		defer ts.Teardown()

		other := ts.Debug().Help.OK.CreateUser("other", "2@tst.tst", "testpass")
		otherClt, _ := ts.Client("2@tst.tst", "testpass")

		otherClt.Help.ERR.EditReaction(
			errors.ErrUnauthorized,
			*reaction.ID, // Someone else's post
			*other.ID,
			"new message",
		)
	})
}
