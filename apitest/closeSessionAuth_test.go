package apitest

import (
	"testing"

	"github.com/romshark/dgraph_graphql_go/api/graph/gqlmod"
	"github.com/romshark/dgraph_graphql_go/apitest/setup"
	"github.com/romshark/dgraph_graphql_go/store/errors"
)

// TestCloseSessionAuth tests session closing authorization
func TestCloseSessionAuth(t *testing.T) {
	setupTest := func(t *testing.T) (
		ts *setup.TestSetup,
		clt *setup.Client,
		session *gqlmod.Session,
	) {
		ts = setup.New(t, tcx)

		ts.Debug().Help.OK.CreateUser(
			"usr",
			"t1@te.te",
			"testpass",
		)
		clt, session = ts.Client("t1@te.te", "testpass")
		return
	}

	// Test closing sessions as a guest
	t.Run("guest (noauth)", func(t *testing.T) {
		ts, _, session := setupTest(t)
		defer ts.Teardown()

		ts.Guest().Help.ERR.CloseSession(
			errors.ErrUnauthorized,
			*session.Key,
		)
	})

	// Test creating closing sessions of other users
	t.Run("non-owner (noauth)", func(t *testing.T) {
		ts, _, session := setupTest(t)
		defer ts.Teardown()

		ts.Debug().Help.OK.CreateUser(
			"other",
			"t2@tst.tst",
			"testpass",
		)
		otherClt, _ := ts.Client("t2@tst.tst", "testpass")

		otherClt.Help.ERR.CloseSession(
			errors.ErrUnauthorized,
			*session.Key,
		)
	})
}
