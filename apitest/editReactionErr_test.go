package apitest

import (
	"testing"

	"github.com/romshark/dgraph_graphql_go/api/graph/gqlmod"
	"github.com/romshark/dgraph_graphql_go/apitest/setup"
	"github.com/romshark/dgraph_graphql_go/store"
	"github.com/romshark/dgraph_graphql_go/store/enum/emotion"
	"github.com/romshark/dgraph_graphql_go/store/errors"
)

// TestEditReactionErr tests all possible reaction editing errors
func TestEditReactionErr(t *testing.T) {
	testSetup := func(t *testing.T) (
		ts *setup.TestSetup,
		debug *setup.Client,
		author *gqlmod.User,
		reaction *gqlmod.Reaction,
	) {
		ts = setup.New(t, tcx)
		debug = ts.Debug()

		author = debug.Help.OK.CreateUser(
			"authoruser",
			"1@tst.tst",
			"testpass",
		)
		post := debug.Help.OK.CreatePost(
			*author.ID,
			"valid title",
			"test contents",
		)
		reaction = debug.Help.OK.CreateReaction(
			*author.ID,
			*post.ID,
			emotion.Thoughtful,
			"test reaction message",
		)
		return
	}

	t.Run("invalidMessage", func(t *testing.T) {
		invalidMessages := map[string]string{
			"empty":   "",
			"tooLong": randomString(257, nil),
		}

		for tName, invalidMessage := range invalidMessages {
			t.Run(tName, func(t *testing.T) {
				ts, debug, author, reaction := testSetup(t)
				defer ts.Teardown()

				debug.Help.ERR.EditReaction(
					errors.ErrInvalidInput,
					*reaction.ID,
					*author.ID,
					invalidMessage,
				)
			})
		}
	})

	t.Run("inexistentReaction", func(t *testing.T) {
		ts, debug, author, _ := testSetup(t)
		defer ts.Teardown()

		debug.Help.ERR.EditReaction(
			errors.ErrInvalidInput,
			store.NewID(), // Inexistent reaction
			*author.ID,
			"new message",
		)
	})

	// Test inexistent editor
	t.Run("inexistentEditor", func(t *testing.T) {
		ts, debug, _, reaction := testSetup(t)
		defer ts.Teardown()

		debug.Help.ERR.EditReaction(
			errors.ErrInvalidInput,
			*reaction.ID,
			store.NewID(), // Inexistent editor
			"new message",
		)
	})
}
