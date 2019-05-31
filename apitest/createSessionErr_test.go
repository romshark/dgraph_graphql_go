package apitest

import (
	"testing"

	"github.com/romshark/dgraph_graphql_go/api/graph/gqlmod"
	"github.com/romshark/dgraph_graphql_go/apitest/setup"
	"github.com/romshark/dgraph_graphql_go/store/errors"
	"github.com/stretchr/testify/require"
)

// TestCreateSessionErr tests all possible sign in errors
func TestCreateSessionErr(t *testing.T) {
	ensureNoSession := func(
		t *testing.T,
		debug *setup.Client,
		user *gqlmod.User,
	) {
		var query struct {
			User *gqlmod.User `json:"user"`
		}
		require.NoError(t, debug.QueryVar(
			`query($userId: Identifier!) {
				user(id: $userId) {
					sessions {
						key
					}
				}
			}`,
			map[string]interface{}{
				"userId": string(*user.ID),
			},
			&query,
		))
		require.NotNil(t, query.User)
		require.Len(t, query.User.Sessions, 0)
	}

	t.Run("wrongEmail", func(t *testing.T) {
		ts := setup.New(t, tcx)
		defer ts.Teardown()

		debug := ts.Debug()

		user := debug.Help.OK.CreateUser(
			"fooBarowich",
			"foo@bar.buz",
			"testpass",
		)
		debug.Help.ERR.CreateSession(
			errors.ErrWrongCreds,
			"foo@fooo.foo",
			"testpass",
		)

		ensureNoSession(t, debug, user)
	})

	t.Run("wrongPassword", func(t *testing.T) {
		ts := setup.New(t, tcx)
		defer ts.Teardown()

		debug := ts.Debug()

		user := debug.Help.OK.CreateUser(
			"fooBarowich",
			"foo@bar.buz",
			"testpass",
		)
		debug.Help.ERR.CreateSession(
			errors.ErrWrongCreds,
			"foo@bar.buz",
			"wronpass",
		)

		ensureNoSession(t, debug, user)
	})

	t.Run("missingEmail", func(t *testing.T) {
		ts := setup.New(t, tcx)
		defer ts.Teardown()

		debug := ts.Debug()

		user := debug.Help.OK.CreateUser(
			"fooBarowich",
			"foo@bar.buz",
			"testpass",
		)
		debug.Help.ERR.CreateSession(errors.ErrInvalidInput, "", "wronpass")

		ensureNoSession(t, debug, user)
	})

	t.Run("missingPassword", func(t *testing.T) {
		ts := setup.New(t, tcx)
		defer ts.Teardown()

		debug := ts.Debug()

		user := debug.Help.OK.CreateUser(
			"fooBarowich",
			"foo@bar.buz",
			"testpass",
		)
		debug.Help.ERR.CreateSession(errors.ErrInvalidInput, "foo@bar.buz", "")

		ensureNoSession(t, debug, user)
	})
}
