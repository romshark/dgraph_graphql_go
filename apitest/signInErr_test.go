package apitest

import (
	"testing"

	"github.com/romshark/dgraph_graphql_go/api/graph/gqlmod"

	"github.com/romshark/dgraph_graphql_go/apitest/setup"
	"github.com/stretchr/testify/require"
)

// TestSignInErr tests all possible sign in errors
func TestSignInErr(t *testing.T) {
	ensureNoSession := func(
		t *testing.T,
		root *setup.Client,
		user *gqlmod.User,
	) {
		var query struct {
			User *gqlmod.User `json:"user"`
		}
		root.QueryVar(
			`query($userId: Identifier!) {
				user(id: $userId) {
					sessions {
						key
					}
				}
			}`,
			map[string]string{
				"userId": string(*user.ID),
			},
			&query,
		)
		require.NotNil(t, query.User)
		require.Len(t, query.User.Sessions, 0)
	}

	t.Run("wrongEmail", func(t *testing.T) {
		ts := setup.New(t, tcx)
		defer ts.Teardown()

		root := ts.Root()

		user := root.Help.OK.CreateUser(
			"fooBarowich",
			"foo@bar.buz",
			"testpass",
		)
		session, err := root.Help.SignIn("foo@fooo.foo", "testpass")
		require.Nil(t, session)
		verifyError(t, "WrongCreds", err)

		ensureNoSession(t, root, user)
	})

	t.Run("wrongPassword", func(t *testing.T) {
		ts := setup.New(t, tcx)
		defer ts.Teardown()

		root := ts.Root()

		user := root.Help.OK.CreateUser(
			"fooBarowich",
			"foo@bar.buz",
			"testpass",
		)
		session, err := root.Help.SignIn("foo@bar.buz", "wronpass")
		require.Nil(t, session)
		verifyError(t, "WrongCreds", err)

		ensureNoSession(t, root, user)
	})

	t.Run("missingEmail", func(t *testing.T) {
		ts := setup.New(t, tcx)
		defer ts.Teardown()

		root := ts.Root()

		user := root.Help.OK.CreateUser(
			"fooBarowich",
			"foo@bar.buz",
			"testpass",
		)
		session, err := root.Help.SignIn("", "wronpass")
		require.Nil(t, session)
		verifyError(t, "InvalidInput", err)

		ensureNoSession(t, root, user)
	})

	t.Run("missingPassword", func(t *testing.T) {
		ts := setup.New(t, tcx)
		defer ts.Teardown()

		root := ts.Root()

		user := root.Help.OK.CreateUser(
			"fooBarowich",
			"foo@bar.buz",
			"testpass",
		)
		session, err := root.Help.SignIn("foo@bar.buz", "")
		require.Nil(t, session)
		verifyError(t, "InvalidInput", err)

		ensureNoSession(t, root, user)
	})
}
