package apitest

import (
	"testing"

	"github.com/romshark/dgraph_graphql_go/api/graph/gqlmod"
	"github.com/romshark/dgraph_graphql_go/apitest/setup"
	"github.com/stretchr/testify/require"
)

// TestCloseAllSessions tests closing all sessions of a user
func TestCloseAllSessions(t *testing.T) {
	t.Run("normal", func(t *testing.T) {
		ts := setup.New(t, tcx)
		defer ts.Teardown()

		// Prepare user
		usrEmail := "usr@tst.tst"
		usrPassword := "testpass"
		usr := ts.Debug().Help.OK.CreateUser("usr", usrEmail, usrPassword)

		// Create 3 sessions
		userClt1, session1 := ts.Client(usrEmail, usrPassword)
		session2 := ts.Debug().Help.OK.CreateSession(usrEmail, usrPassword)
		session3 := ts.Debug().Help.OK.CreateSession(usrEmail, usrPassword)

		expectedKeys := []string{
			*session1.Key,
			*session2.Key,
			*session3.Key,
		}

		// Close all sessions
		keys := userClt1.Help.OK.CloseAllSessions(*usr.ID)

		require.Equal(t, expectedKeys, keys)

		// Query all sessions of the user and assume none to be found
		var query struct {
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
			&query,
		))
		require.NotNil(t, query.User)
		require.Len(t, query.User.Sessions, 0)
	})

	t.Run("noSessions", func(t *testing.T) {
		ts := setup.New(t, tcx)
		defer ts.Teardown()

		// Prepare user
		usr := ts.Debug().Help.OK.CreateUser("usr", "usr@tst.tst", "testpass")

		// Close all sessions
		keys := ts.Debug().Help.OK.CloseAllSessions(*usr.ID)
		require.Len(t, keys, 0)

		// Query all sessions of the user and assume none to be found
		var query struct {
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
			&query,
		))
		require.NotNil(t, query.User)
		require.Len(t, query.User.Sessions, 0)
	})
}
