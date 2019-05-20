package apitest

import (
	"testing"

	"github.com/romshark/dgraph_graphql_go/apitest/setup"
	"github.com/romshark/dgraph_graphql_go/store"
	strerr "github.com/romshark/dgraph_graphql_go/store/errors"
)

// TestCloseAllSessionsErr tests all possible errors for all session closing
func TestCloseAllSessionsErr(t *testing.T) {
	t.Run("inexistentUser", func(t *testing.T) {
		ts := setup.New(t, tcx)
		defer ts.Teardown()

		ts.Debug().Help.ERR.CloseAllSessions(
			strerr.ErrInvalidInput,
			store.NewID(),
		)
	})
}
