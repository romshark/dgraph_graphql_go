package apitest

import (
	"testing"

	"github.com/romshark/dgraph_graphql_go/apitest/setup"
	strerr "github.com/romshark/dgraph_graphql_go/store/errors"
)

// TestCloseSessionErr tests all possible sessions closing errors
func TestCloseSessionErr(t *testing.T) {
	t.Run("inexistentSession", func(t *testing.T) {
		ts := setup.New(t, tcx)
		defer ts.Teardown()

		ts.Debug().Help.ERR.CloseSession(
			strerr.ErrInvalidInput,
			"inexistent_session_key",
		)
	})

	t.Run("repeatedClosing", func(t *testing.T) {
		ts := setup.New(t, tcx)
		defer ts.Teardown()

		// Prepare user
		usrEmail := "2@tst.tst"
		usrPassword := "testpass"
		ts.Debug().Help.OK.CreateUser("usr", usrEmail, usrPassword)

		// Create user session
		_, session1 := ts.Client(usrEmail, usrPassword)

		ts.Debug().Help.OK.CloseSession(*session1.Key)
		ts.Debug().Help.ERR.CloseSession(
			strerr.ErrInvalidInput,
			"inexistent_session_key",
		)
	})
}
