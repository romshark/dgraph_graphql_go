package apitest

import (
	"testing"

	"github.com/romshark/dgraph_graphql_go/apitest/setup"
	"github.com/romshark/dgraph_graphql_go/store"
	"github.com/romshark/dgraph_graphql_go/store/errors"
)

// TestCreatePostErr tests all possible post creation errors
func TestCreatePostErr(t *testing.T) {
	// Test invalid title
	t.Run("invalidTitle", func(t *testing.T) {
		invalidTitles := map[string]string{
			"empty":    "",
			"tooShort": "f",
			"tooLong":  randomString(65, nil),
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
				debug.Help.ERR.CreatePost(
					errors.ErrInvalidInput,
					*author.ID,
					invalidTitle,
					"test contents",
				)
			})
		}
	})

	// Test invalid contents
	t.Run("invalidContents", func(t *testing.T) {
		invalidContents := map[string]string{
			"empty":   "",
			"tooLong": randomString(257, nil),
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
				debug.Help.ERR.CreatePost(
					errors.ErrInvalidInput,
					*author.ID,
					"test title",
					invalidContent,
				)
			})
		}
	})

	// Test inexistent author
	t.Run("inexistentAuthor", func(t *testing.T) {
		ts := setup.New(t, tcx)
		defer ts.Teardown()

		ts.Debug().Help.ERR.CreatePost(
			errors.ErrInvalidInput,
			store.NewID(),
			"test title",
			"test contents",
		)
	})
}
