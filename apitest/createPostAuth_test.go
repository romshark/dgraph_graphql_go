package apitest

import (
	"testing"

	"github.com/romshark/dgraph_graphql_go/api/graph/gqlmod"
	"github.com/romshark/dgraph_graphql_go/apitest/setup"
	"github.com/stretchr/testify/require"
)

// TestCreatePostAuth tests post creation authorization
func TestCreatePostAuth(t *testing.T) {
	setupTest := func(t *testing.T) (
		ts *setup.TestSetup,
		author *gqlmod.User,
		authorClt *setup.Client,
	) {
		ts = setup.New(t, tcx)
		root := ts.Root()

		authorEmail := "t1@te.te"
		authorPass := "testpass"
		author = root.Help.OK.CreateUser(
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

		post, err := ts.Guest().Help.CreatePost(
			*author.ID,
			"test post",
			"test content",
		)
		require.Nil(t, post)
		verifyError(t, "Unauthorized", err)
	})

	// Test creating posts on behalf of other users
	t.Run("non-author (noauth)", func(t *testing.T) {
		ts, _, clt := setupTest(t)
		defer ts.Teardown()

		other := ts.Root().Help.OK.CreateUser("other", "t2@tst.tst", "testpass")
		post, err := clt.Help.CreatePost(*other.ID, "test post", "test content")
		require.Nil(t, post)
		verifyError(t, "Unauthorized", err)
	})
}
