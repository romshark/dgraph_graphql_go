package apitest

import (
	"testing"

	"github.com/romshark/dgraph_graphql_go/apitest/setup"
	"github.com/stretchr/testify/require"
)

// TestCreateUserErr tests all possible user account creation errors
func TestCreateUserErr(t *testing.T) {
	// Test reserved email on creation
	t.Run("reservedEmail", func(t *testing.T) {
		ts := setup.New(t, tcx)
		defer ts.Teardown()

		clt := ts.Root()

		clt.Help.OK.CreateUser("fooBarowich", "foo@bar.buz", "testpass")
		res, err := clt.Help.CreateUser(
			"bazBuzowich",
			"foo@bar.buz",
			"testpass",
		)
		require.Nil(t, res)
		verifyError(t, "InvalidInput", err)
	})

	// Test reserved displayName on creation
	t.Run("reservedDisplayName", func(t *testing.T) {
		ts := setup.New(t, tcx)
		defer ts.Teardown()

		clt := ts.Root()

		clt.Help.OK.CreateUser("fooBarowich", "foo@bar.buz", "testpass")
		res, err := clt.Help.CreateUser(
			"fooBarowich",
			"baz@buzowich.buz",
			"testpass",
		)
		require.Nil(t, res)
		verifyError(t, "InvalidInput", err)
	})

	// Test reserved displayName on creation
	t.Run("invalidDisplayName", func(t *testing.T) {
		invalidDisplayNames := map[string]string{
			"empty":    "",
			"tooShort": "t",
			"tooLong": "11111111000000001111111100000000" +
				"11111111000000001111111100000000" +
				"2",
		}

		for tName, invalidDisplayName := range invalidDisplayNames {
			t.Run(tName, func(t *testing.T) {
				ts := setup.New(t, tcx)
				defer ts.Teardown()

				clt := ts.Root()

				res, err := clt.Help.CreateUser(
					invalidDisplayName,
					"test@test.test",
					"foobar",
				)
				require.Nil(t, res)
				verifyError(t, "InvalidInput", err)
			})
		}
	})

	// Test reserved displayName on creation
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

				clt := ts.Root()

				res, err := clt.Help.CreateUser(
					"testDisplayName",
					invalidEmail,
					"testpass",
				)
				require.Nil(t, res)
				verifyError(t, "InvalidInput", err)
			})
		}
	})
}
