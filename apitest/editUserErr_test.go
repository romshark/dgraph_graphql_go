package apitest

import (
	"testing"

	"github.com/romshark/dgraph_graphql_go/api/graph/gqlmod"
	"github.com/romshark/dgraph_graphql_go/apitest/setup"
	"github.com/romshark/dgraph_graphql_go/store"
	"github.com/romshark/dgraph_graphql_go/store/errors"
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
			"tooLong": randomString(97, nil),
		}

		for tName, invalidEmail := range invalidEmails {
			t.Run(tName, func(t *testing.T) {
				ts, debug, user := testSetup(t)
				defer ts.Teardown()

				debug.Help.ERR.EditUser(
					errors.ErrInvalidInput,
					*user.ID,
					*user.ID,
					&invalidEmail,
					nil,
				)
			})
		}
	})

	t.Run("invalidPasswords", func(t *testing.T) {
		invalidPasswords := map[string]string{
			"empty":    "",
			"tooShort": "12345",
			"tooLong":  randomString(257, nil),
		}

		for tName, invalidPassword := range invalidPasswords {
			t.Run(tName, func(t *testing.T) {
				ts, debug, user := testSetup(t)
				defer ts.Teardown()

				debug.Help.ERR.EditUser(
					errors.ErrInvalidInput,
					*user.ID,
					*user.ID,
					nil,
					&invalidPassword,
				)
			})
		}
	})

	t.Run("inexistentUser", func(t *testing.T) {
		ts, debug, user := testSetup(t)
		defer ts.Teardown()

		newEmail := "new@email.test"
		debug.Help.ERR.EditUser(
			errors.ErrInvalidInput,
			store.NewID(), // Inexistent profile
			*user.ID,
			&newEmail,
			nil,
		)
	})

	t.Run("inexistentEditor", func(t *testing.T) {
		ts, debug, user := testSetup(t)
		defer ts.Teardown()

		newEmail := "new@email.test"
		debug.Help.ERR.EditUser(
			errors.ErrInvalidInput,
			*user.ID,
			store.NewID(), // Inexistent editor
			&newEmail,
			nil,
		)
	})

	t.Run("noChanges", func(t *testing.T) {
		ts, debug, user := testSetup(t)
		defer ts.Teardown()

		debug.Help.ERR.EditUser(
			errors.ErrInvalidInput,
			*user.ID,
			*user.ID,
			nil,
			nil,
		)
	})
}
