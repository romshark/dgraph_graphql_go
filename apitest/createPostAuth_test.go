package apitest

import (
	"testing"

	"github.com/romshark/dgraph_graphql_go/api/graph/gqlmod"
	"github.com/romshark/dgraph_graphql_go/apitest/setup"
	"github.com/romshark/dgraph_graphql_go/store/errors"
)

// TestCreatePostAuth tests post creation authorization
func TestCreatePostAuth(t *testing.T) {
	setupTest := func(t *testing.T) (
		ts *setup.TestSetup,
		author *gqlmod.User,
		authorClt *setup.Client,
	) {
		ts = setup.New(t, tcx)
		debug := ts.Debug()

		authorEmail := "t1@te.te"
		authorPass := "testpass"
		author = debug.Help.OK.CreateUser(
			"fooBarowich",
			authorEmail,
			authorPass,
		)
		authorClt, _ = ts.Client(authorEmail, authorPass)
		return
	}

	// Test creating posts as a guest
	t.Run("guest (noauth)", func(t *testing.T) {
		ts, author, _ := setupTest(t)
		defer ts.Teardown()

		ts.Guest().Help.ERR.CreatePost(
			errors.ErrUnauthorized,
			*author.ID,
			"test post",
			"test content",
		)
	})

	// Test creating posts on behalf of other users
	t.Run("non-author (noauth)", func(t *testing.T) {
		ts, _, clt := setupTest(t)
		defer ts.Teardown()

		other := ts.Debug().Help.OK.CreateUser(
			"other",
			"t2@tst.tst",
			"testpass",
		)
		clt.Help.ERR.CreatePost(
			errors.ErrUnauthorized,
			*other.ID,
			"test post",
			"test content",
		)
	})
}
