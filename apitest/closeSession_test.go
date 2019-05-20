package apitest

import (
	"testing"

	"github.com/romshark/dgraph_graphql_go/api/graph/gqlmod"
	"github.com/romshark/dgraph_graphql_go/apitest/setup"
	"github.com/stretchr/testify/require"
)

// TestCloseSession tests closing sessions by key
func TestCloseSession(t *testing.T) {
	ts := setup.New(t, tcx)
	defer ts.Teardown()

	// Prepare user
	usrEmail := "2@tst.tst"
	usrPassword := "testpass"
	usr := ts.Debug().Help.OK.CreateUser("usr", usrEmail, usrPassword)

	// Create user sessions
	usrClt1, session1 := ts.Client(usrEmail, usrPassword)
	usrClt2, session2 := ts.Client(usrEmail, usrPassword)

	// Close first session of the user
	usrClt1.Help.OK.CloseSession(*session1.Key)

	// Query all sessions of the user and assume the second to be found
	var queryAfter1 struct {
		User *gqlmod.User `json:"user"`
	}
	require.NoError(t, ts.Debug().QueryVar(
		`query($userId: Identifier!) {
			user(id: $userId) {
				sessions {
					key
					creation
				}
			}
		}`,
		map[string]interface{}{
			"userId": string(*usr.ID),
		},
		&queryAfter1,
	))
	require.NotNil(t, queryAfter1.User)
	require.Len(t, queryAfter1.User.Sessions, 1)
	require.Equal(t, queryAfter1.User.Sessions[0].Key, session2.Key)
	require.Equal(t, queryAfter1.User.Sessions[0].Creation, session2.Creation)

	// Close second session of the user
	usrClt2.Help.OK.CloseSession(*session2.Key)

	// Query all sessions of the user and assume none to be found
	var queryAfter2 struct {
		User *gqlmod.User `json:"user"`
	}
	require.NoError(t, ts.Debug().QueryVar(
		`query($userId: Identifier!) {
			user(id: $userId) {
				sessions {
					key
					creation
				}
			}
		}`,
		map[string]interface{}{
			"userId": string(*usr.ID),
		},
		&queryAfter2,
	))
	require.NotNil(t, queryAfter2.User)
	require.Len(t, queryAfter2.User.Sessions, 0)
}
