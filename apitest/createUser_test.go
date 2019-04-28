package apitest

import (
	"testing"
	"time"

	"github.com/romshark/dgraph_graphql_go/api/graph/gqlmod"
	"github.com/romshark/dgraph_graphql_go/apitest/setup"
	"github.com/stretchr/testify/require"
)

// TestCreateCustomer tests customer account creation
func TestCreateCustomer(t *testing.T) {
	ts := setup.New(t, tcx)
	defer ts.Teardown()

	var result struct {
		CreateUser *gqlmod.User `json:"createUser"`
	}
	ts.Query(
		`mutation {
			createUser(
				email: "foo@bar.buz"
				displayName: "fooBar",
			) {
				id
				email
				displayName
				creation
				posts {
					id
				}
			}
		}`,
		&result,
	)

	require.NotNil(t, result.CreateUser)
	require.Len(t, *result.CreateUser.ID, 32)
	require.Equal(t, "foo@bar.buz", *result.CreateUser.Email)
	require.Equal(t, "fooBar", *result.CreateUser.DisplayName)
	require.Len(t, result.CreateUser.Posts, 0)
	require.WithinDuration(
		t,
		time.Now(),
		*result.CreateUser.Creation,
		time.Second*3,
	)
}
