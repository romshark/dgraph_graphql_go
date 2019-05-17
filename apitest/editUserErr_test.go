package apitest

import (
	"testing"

	"github.com/romshark/dgraph_graphql_go/api/graph/gqlmod"
	"github.com/romshark/dgraph_graphql_go/apitest/setup"
	"github.com/romshark/dgraph_graphql_go/store"
	"github.com/stretchr/testify/require"
)

// TestEditUserErr tests all possible user profile editing errors
func TestEditUserErr(t *testing.T) {
	testSetup := func(t *testing.T) (
		ts *setup.TestSetup,
		debug *setup.Client,
		user *gqlmod.User,
	) {
		ts = setup.New(t, tcx)
		debug = ts.Debug()

		user = debug.Help.OK.CreateUser(
			"testuser",
			"test@test.test",
			"testpass",
		)
		return
	}

	t.Run("invalidEmail", func(t *testing.T) {
		invalidEmails := map[string]string{
			"empty":   "",
			"invalid": "fooooo@baaaaaar",
			"tooLong": "11111111000000001111111100000000" +
				"11111111000000001111111100000000" +
				"11111111000000001111111100000000" +
				"2",
		}

		for tName, invalidEmail := range invalidEmails {
			t.Run(tName, func(t *testing.T) {
				ts, debug, user := testSetup(t)
				defer ts.Teardown()

				res, err := debug.Help.EditUser(
					*user.ID,
					*user.ID,
					&invalidEmail,
					nil,
				)
				require.Nil(t, res)
				verifyError(t, "InvalidInput", err)
			})
		}
	})

	t.Run("invalidPasswords", func(t *testing.T) {
		invalidPasswords := map[string]string{
			"empty":    "",
			"tooShort": "12345",
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

		for tName, invalidPassword := range invalidPasswords {
			t.Run(tName, func(t *testing.T) {
				ts, debug, user := testSetup(t)
				defer ts.Teardown()

				res, err := debug.Help.EditUser(
					*user.ID,
					*user.ID,
					nil,
					&invalidPassword,
				)
				require.Nil(t, res)
				verifyError(t, "InvalidInput", err)
			})
		}
	})

	t.Run("inexistentUser", func(t *testing.T) {
		ts, debug, user := testSetup(t)
		defer ts.Teardown()

		newEmail := "new@email.test"
		res, err := debug.Help.EditUser(
			store.NewID(), // Inexistent profile
			*user.ID,
			&newEmail,
			nil,
		)
		require.Nil(t, res)
		verifyError(t, "InvalidInput", err)
	})

	t.Run("inexistentEditor", func(t *testing.T) {
		ts, debug, user := testSetup(t)
		defer ts.Teardown()

		newEmail := "new@email.test"
		res, err := debug.Help.EditUser(
			*user.ID,
			store.NewID(), // Inexistent editor
			&newEmail,
			nil,
		)
		require.Nil(t, res)
		verifyError(t, "InvalidInput", err)
	})

	t.Run("noChanges", func(t *testing.T) {
		ts, debug, user := testSetup(t)
		defer ts.Teardown()

		res, err := debug.Help.EditUser(
			*user.ID,
			*user.ID,
			nil,
			nil,
		)
		require.Nil(t, res)
		verifyError(t, "InvalidInput", err)
	})
}
