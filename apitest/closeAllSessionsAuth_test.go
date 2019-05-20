package apitest

import (
	"testing"

	"github.com/romshark/dgraph_graphql_go/api/graph/gqlmod"
	"github.com/romshark/dgraph_graphql_go/apitest/setup"
	"github.com/romshark/dgraph_graphql_go/store/errors"
)

// TestCloseAllSessionsAuth tests all sessions closing authorization
func TestCloseAllSessionsAuth(t *testing.T) {
	setupTest := func(t *testing.T) (
		ts *setup.TestSetup,
		usr *gqlmod.User,
		sessions []*gqlmod.Session,
	) {
		ts = setup.New(t, tcx)

		usr = ts.Debug().Help.OK.CreateUser(
			"usr",
			"t1@te.te",
			"testpass",
		)
		_, session1 := ts.Client("t1@te.te", "testpass")
		_, session2 := ts.Client("t1@te.te", "testpass")
		sessions = []*gqlmod.Session{session1, session2}
		return
	}

	// Test closing all sessions of a user as a guest
	t.Run("guest (noauth)", func(t *testing.T) {
		ts, usr, _ := setupTest(t)
		defer ts.Teardown()

		ts.Guest().Help.ERR.CloseAllSessions(
			errors.ErrUnauthorized,
			*usr.ID,
		)
	})

	// Test creating closing all sessions of other users
	t.Run("non-owner (noauth)", func(t *testing.T) {
		ts, usr, _ := setupTest(t)
		defer ts.Teardown()

		ts.Debug().Help.OK.CreateUser(
			"other",
			"t2@tst.tst",
			"testpass",
		)
		otherClt, _ := ts.Client("t2@tst.tst", "testpass")

		otherClt.Help.ERR.CloseAllSessions(
			errors.ErrUnauthorized,
			*usr.ID,
		)
	})
}
