package apitest

import (
	"testing"

	"github.com/romshark/dgraph_graphql_go/apitest/setup"
	"github.com/romshark/dgraph_graphql_go/store"
	"github.com/romshark/dgraph_graphql_go/store/enum/emotion"
	"github.com/stretchr/testify/require"
)

// TestCreateReactionErr tests all possible reaction creation errors
func TestCreateReactionErr(t *testing.T) {
	t.Run("inexistentAuthor", func(t *testing.T) {
		ts := setup.New(t, tcx)
		defer ts.Teardown()

		debug := ts.Debug()

		// User 1
		firstP := debug.Help.OK.CreateUser("first", "1@test.test", "testpass")
		post := debug.Help.OK.CreatePost(*firstP.ID, "Test", "test")

		res, err := debug.Help.CreateReaction(
			store.NewID(), // inexistent author
			*post.ID,
			emotion.Excited,
			"test message",
		)
		require.Nil(t, res)
		verifyError(t, "InvalidInput", err)
	})

	t.Run("inexistentPost", func(t *testing.T) {
		ts := setup.New(t, tcx)
		defer ts.Teardown()

		debug := ts.Debug()

		// User 1
		firstP := debug.Help.OK.CreateUser("first", "1@test.test", "testpass")

		res, err := debug.Help.CreateReaction(
			*firstP.ID,
			store.NewID(), // inexistent post
			emotion.Excited,
			"test message",
		)
		require.Nil(t, res)
		verifyError(t, "InvalidInput", err)
	})

	t.Run("invalidMessage", func(t *testing.T) {
		invalidMessages := map[string]string{
			"empty": "",
			"tooLong": "11111111000000001111111100000000" +
				"11111111000000001111111100000000" +
				"11111111000000001111111100000000" +
				"11111111000000001111111100000000" +
				"11111111000000001111111100000000" +
				"11111111000000001111111100000000" +
				"11111111000000001111111100000000" +
				"11111111000000001111111100000000" +
				"2",
		}

		for caseName, invalidMessage := range invalidMessages {
			t.Run(caseName, func(t *testing.T) {
				ts := setup.New(t, tcx)
				defer ts.Teardown()

				debug := ts.Debug()

				// User 1
				firstP := debug.Help.OK.CreateUser(
					"first",
					"1@test.test",
					"testpass",
				)
				post := debug.Help.OK.CreatePost(*firstP.ID, "Test", "test")
				secondP := debug.Help.OK.CreateUser(
					"second",
					"2@test.test",
					"testpass",
				)

				res, err := debug.Help.CreateReaction(
					*secondP.ID,
					*post.ID,
					emotion.Excited,
					invalidMessage,
				)
				require.Nil(t, res)
				verifyError(t, "InvalidInput", err)
			})
		}
	})
}
