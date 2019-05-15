package apitest

import (
	"testing"

	"github.com/romshark/dgraph_graphql_go/apitest/setup"
	"github.com/romshark/dgraph_graphql_go/store"
	"github.com/stretchr/testify/require"
)

// TestCreatePostErr tests all possible post creation errors
func TestCreatePostErr(t *testing.T) {
	// Test invalid title
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
				ts := setup.New(t, tcx)
				defer ts.Teardown()

				debug := ts.Debug()

				author := debug.Help.OK.CreateUser(
					"fooBarowich",
					"foo@bar.buz",
					"testpass",
				)
				res, err := debug.Help.CreatePost(
					*author.ID,
					invalidTitle,
					"test contents",
				)
				require.Nil(t, res)
				verifyError(t, "InvalidInput", err)
			})
		}
	})

	// Test invalid contents
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
				ts := setup.New(t, tcx)
				defer ts.Teardown()

				debug := ts.Debug()

				author := debug.Help.OK.CreateUser(
					"fooBarowich",
					"foo@bar.buz",
					"testpass",
				)
				res, err := debug.Help.CreatePost(
					*author.ID,
					"test title",
					invalidContent,
				)
				require.Nil(t, res)
				verifyError(t, "InvalidInput", err)
			})
		}
	})

	// Test inexistent author
	t.Run("inexistentAuthor", func(t *testing.T) {
		ts := setup.New(t, tcx)
		defer ts.Teardown()

		res, err := ts.Debug().Help.CreatePost(
			store.NewID(),
			"test title",
			"test contents",
		)
		require.Nil(t, res)
		verifyError(t, "InvalidInput", err)
	})
}
