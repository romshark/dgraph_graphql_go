package apitest

import (
	"testing"

	"github.com/romshark/dgraph_graphql_go/apitest/setup"
	"github.com/stretchr/testify/require"
)

// TestEditUser tests user profile editing
func TestEditUser(t *testing.T) {
	t.Run("email", func(t *testing.T) {
		ts := setup.New(t, tcx)
		defer ts.Teardown()

		// Prepare
		oldEmail := "1@tst.tst"
		password := "testpass"
		debug := ts.Debug()
		user := debug.Help.OK.CreateUser("user", oldEmail, password)
		userClt, _ := ts.Client(oldEmail, password)

		// Test edit
		newEmail := "new@email.test"
		userClt.Help.OK.EditUser(
			*user.ID,
			*user.ID,
			&newEmail,
			nil, // don't change the password
		)

		// Test signing in using the old email
		session, err := ts.Guest().Help.SignIn(oldEmail, password)
		require.Nil(t, session)
		verifyError(t, "WrongCreds", err)

		// Test signing in using the new email
		ts.Guest().Help.OK.SignIn(newEmail, password)
	})

	t.Run("password", func(t *testing.T) {
		ts := setup.New(t, tcx)
		defer ts.Teardown()

		email := "1@tst.tst"
		oldPassword := "testpass"

		// Prepare
		debug := ts.Debug()
		user := debug.Help.OK.CreateUser("user", email, oldPassword)
		userClt, _ := ts.Client(email, oldPassword)

		// Test edit
		newPassword := "newpassword"
		userClt.Help.OK.EditUser(
			*user.ID,
			*user.ID,
			nil, // don't change the email
			&newPassword,
		)

		// Test signing in using the old password
		session, err := ts.Guest().Help.SignIn(email, oldPassword)
		require.Nil(t, session)
		verifyError(t, "WrongCreds", err)

		// Test signing in using the new password
		ts.Guest().Help.OK.SignIn(email, newPassword)
	})
}
