package apitest

import (
	"testing"

	"github.com/romshark/dgraph_graphql_go/apitest/setup"
	"github.com/romshark/dgraph_graphql_go/store/errors"
)

// TestCreateUserErr tests all possible user account creation errors
func TestCreateUserErr(t *testing.T) {
	// Test reserved email on creation
	t.Run("reservedEmail", func(t *testing.T) {
		ts := setup.New(t, tcx)
		defer ts.Teardown()

		debug := ts.Debug()

		debug.Help.OK.CreateUser("fooBarowich", "foo@bar.buz", "testpass")
		debug.Help.ERR.CreateUser(
			errors.ErrInvalidInput,
			"bazBuzowich",
			"foo@bar.buz",
			"testpass",
		)
	})

	// Test reserved displayName on creation
	t.Run("reservedDisplayName", func(t *testing.T) {
		ts := setup.New(t, tcx)
		defer ts.Teardown()

		debug := ts.Debug()

		debug.Help.OK.CreateUser("fooBarowich", "foo@bar.buz", "testpass")
		debug.Help.ERR.CreateUser(
			errors.ErrInvalidInput,
			"fooBarowich",
			"baz@buzowich.buz",
			"testpass",
		)
	})

	// Test invalid displayName on creation
	t.Run("invalidDisplayName", func(t *testing.T) {
		invalidDisplayNames := map[string]string{
			"empty":    "",
			"tooShort": "t",
			"tooLong":  randomString(65, nil),
		}

		for tName, invalidDisplayName := range invalidDisplayNames {
			t.Run(tName, func(t *testing.T) {
				ts := setup.New(t, tcx)
				defer ts.Teardown()

				ts.Debug().Help.ERR.CreateUser(
					errors.ErrInvalidInput,
					invalidDisplayName,
					"test@test.test",
					"foobar",
				)
			})
		}
	})

	// Test invalid email on creation
	t.Run("invalidEmail", func(t *testing.T) {
		invalidEmails := map[string]string{
			"empty":       "",
			"missingTld":  "test@test",
			"missingHost": "test",
			"tooLong": "teeeeeeeeeeeeeeeeeeeeeeeeeeeeeee" +
				"eeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee" +
				"eeeeeeeeeeeeeeeeeeee@teeest.test" +
				"t",
		}

		for tName, invalidEmail := range invalidEmails {
			t.Run(tName, func(t *testing.T) {
				ts := setup.New(t, tcx)
				defer ts.Teardown()

				ts.Debug().Help.ERR.CreateUser(
					errors.ErrInvalidInput,
					"testDisplayName",
					invalidEmail,
					"testpass",
				)
			})
		}
	})
}
