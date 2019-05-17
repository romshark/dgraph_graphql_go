package apitest

import (
	"testing"

	"github.com/romshark/dgraph_graphql_go/api/graph/gqlmod"
	"github.com/romshark/dgraph_graphql_go/apitest/setup"
	"github.com/romshark/dgraph_graphql_go/store/errors"
)

// TestEditUserAuth tests user profile editing authorization
func TestEditUserAuth(t *testing.T) {
	setupTest := func(t *testing.T) (
		ts *setup.TestSetup,
		user *gqlmod.User,
		userClt *setup.Client,
	) {
		ts = setup.New(t, tcx)
		debug := ts.Debug()

		userEmail := "user@tst.tst"
		userPass := "testpass"
		user = debug.Help.OK.CreateUser(
			"testuser",
			userEmail,
			userPass,
		)
		userClt, _ = ts.Client(userEmail, userPass)
		return
	}

	// Test editing profiles as a guest
	t.Run("guest (noauth)", func(t *testing.T) {
		ts, user, _ := setupTest(t)
		defer ts.Teardown()

		newEmail := "new@email.test"
		newPassword := "newpassword"
		ts.Guest().Help.ERR.EditUser(
			errors.ErrUnauthorized,
			*user.ID,
			*user.ID,
			&newEmail,
			&newPassword,
		)
	})

	// Test editing profiles on behalf of other users
	t.Run("non-editor (noauth)", func(t *testing.T) {
		ts, user, userClt := setupTest(t)
		defer ts.Teardown()

		other := ts.Debug().Help.OK.CreateUser("other", "2@tst.tst", "testpass")

		newEmail := "new@email.test"
		newPassword := "newpassword"
		userClt.Help.ERR.EditUser(
			errors.ErrUnauthorized,
			*user.ID,
			*other.ID, // Different editor ID
			&newEmail,
			&newPassword,
		)
	})

	// Test editing profiles of other users
	t.Run("non-owner (noauth)", func(t *testing.T) {
		ts, user, _ := setupTest(t)
		defer ts.Teardown()

		other := ts.Debug().Help.OK.CreateUser("other", "2@tst.tst", "testpass")
		otherClt, _ := ts.Client("2@tst.tst", "testpass")

		newEmail := "new@email.test"
		newPassword := "newpassword"
		otherClt.Help.ERR.EditUser(
			errors.ErrUnauthorized,
			*user.ID, // Someone else's profile
			*other.ID,
			&newEmail,
			&newPassword,
		)
	})
}
