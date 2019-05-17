package apitest

import (
	"testing"

	"github.com/romshark/dgraph_graphql_go/api/graph/gqlmod"
	"github.com/romshark/dgraph_graphql_go/apitest/setup"
	"github.com/romshark/dgraph_graphql_go/store"
	"github.com/romshark/dgraph_graphql_go/store/errors"
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

				debug.Help.ERR.EditPost(
					errors.ErrInvalidInput,
					*post.ID,
					*author.ID,
					&invalidTitle,
					post.Contents,
				)
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

				debug.Help.ERR.EditPost(
					errors.ErrInvalidInput,
					*post.ID,
					*author.ID,
					post.Title,
					&invalidContent,
				)
			})
		}
	})

	t.Run("inexistentPost", func(t *testing.T) {
		ts, debug, post, _ := testSetup(t)
		defer ts.Teardown()

		debug.Help.ERR.EditPost(
			errors.ErrInvalidInput,
			store.NewID(), // Inexistent post
			*post.ID,
			post.Title,
			post.Contents,
		)
	})

	// Test inexistent editor
	t.Run("inexistentEditor", func(t *testing.T) {
		ts, debug, post, _ := testSetup(t)
		defer ts.Teardown()

		debug.Help.ERR.EditPost(
			errors.ErrInvalidInput,
			*post.ID,
			store.NewID(), // Inexistent editor
			post.Title,
			post.Contents,
		)
	})

	t.Run("noChanges", func(t *testing.T) {
		ts, debug, post, _ := testSetup(t)
		defer ts.Teardown()

		debug.Help.ERR.EditPost(
			errors.ErrInvalidInput,
			*post.ID,
			*post.ID,
			nil,
			nil,
		)
	})
}
