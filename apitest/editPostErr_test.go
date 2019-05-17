package apitest

import (
	"testing"

	"github.com/romshark/dgraph_graphql_go/api/graph/gqlmod"
	"github.com/romshark/dgraph_graphql_go/apitest/setup"
	"github.com/romshark/dgraph_graphql_go/store"
	"github.com/stretchr/testify/require"
)

// TestEditPostErr tests all possible post editing errors
func TestEditPostErr(t *testing.T) {
	testSetup := func(t *testing.T) (
		ts *setup.TestSetup,
		debug *setup.Client,
		post *gqlmod.Post,
		author *gqlmod.User,
	) {
		ts = setup.New(t, tcx)
		debug = ts.Debug()

		author = debug.Help.OK.CreateUser(
			"fooBarowich",
			"foo@bar.buz",
			"testpass",
		)
		post = debug.Help.OK.CreatePost(
			*author.ID,
			"valid title",
			"test contents",
		)
		return
	}

	t.Run("invalidTitle", func(t *testing.T) {
		invalidTitles := map[string]string{
			"empty":    "",
			"tooShort": "f",
			"tooLong": "11111111000000001111111100000000" +
				"11111111000000001111111100000000" +
				"2",
		}

		for tName, invalidTitle := range invalidTitles {
			t.Run(tName, func(t *testing.T) {
				ts, debug, post, author := testSetup(t)
				defer ts.Teardown()

				res, err := debug.Help.EditPost(
					*post.ID,
					*author.ID,
					&invalidTitle,
					post.Contents,
				)
				require.Nil(t, res)
				verifyError(t, "InvalidInput", err)
			})
		}
	})

	t.Run("invalidContents", func(t *testing.T) {
		invalidContents := map[string]string{
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

		for tName, invalidContent := range invalidContents {
			t.Run(tName, func(t *testing.T) {
				ts, debug, post, author := testSetup(t)
				defer ts.Teardown()

				res, err := debug.Help.EditPost(
					*post.ID,
					*author.ID,
					post.Title,
					&invalidContent,
				)
				require.Nil(t, res)
				verifyError(t, "InvalidInput", err)
			})
		}
	})

	t.Run("inexistentPost", func(t *testing.T) {
		ts, debug, post, _ := testSetup(t)
		defer ts.Teardown()

		res, err := debug.Help.EditPost(
			store.NewID(), // Inexistent post
			*post.ID,
			post.Title,
			post.Contents,
		)
		require.Nil(t, res)
		verifyError(t, "InvalidInput", err)
	})

	// Test inexistent editor
	t.Run("inexistentEditor", func(t *testing.T) {
		ts, debug, post, _ := testSetup(t)
		defer ts.Teardown()

		res, err := debug.Help.EditPost(
			*post.ID,
			store.NewID(), // Inexistent editor
			post.Title,
			post.Contents,
		)
		require.Nil(t, res)
		verifyError(t, "InvalidInput", err)
	})

	t.Run("noChanges", func(t *testing.T) {
		ts, debug, post, _ := testSetup(t)
		defer ts.Teardown()

		res, err := debug.Help.EditPost(
			*post.ID,
			*post.ID,
			nil,
			nil,
		)
		require.Nil(t, res)
		verifyError(t, "InvalidInput", err)
	})
}
