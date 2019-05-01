package apitest

import (
	"testing"

	"github.com/romshark/dgraph_graphql_go/api"
	"github.com/romshark/dgraph_graphql_go/api/graph/gqlmod"
	"github.com/stretchr/testify/require"
)

// verifyError helps verifying an API error
func verifyError(
	t *testing.T,
	expectedCode string,
	err *api.ResponseError,
) {
	require.NotNil(t, err)
	require.Equal(t, expectedCode, err.Code)
	require.True(t, len(err.Message) > 0)
}

func compareUsers(
	t *testing.T,
	expected *gqlmod.User,
	actual *gqlmod.User,
) {
	require.NotNil(t, expected)
	require.NotNil(t, actual)

	// id
	require.NotNil(t, expected.ID)
	require.NotNil(t, actual.ID)
	require.Equal(t, *expected.ID, *actual.ID)

	// email
	require.NotNil(t, expected.Email)
	require.NotNil(t, actual.Email)
	require.Equal(t, *expected.Email, *actual.Email)

	// displayName
	require.NotNil(t, expected.DisplayName)
	require.NotNil(t, actual.DisplayName)
	require.Equal(t, *expected.DisplayName, *actual.DisplayName)

	// creation
	require.NotNil(t, expected.Creation)
	require.NotNil(t, actual.Creation)
	require.Equal(t, *expected.Creation, *actual.Creation)
}

func comparePosts(
	t *testing.T,
	expected *gqlmod.Post,
	actual *gqlmod.Post,
) {
	require.NotNil(t, expected)
	require.NotNil(t, actual)

	// id
	require.NotNil(t, expected.ID)
	require.NotNil(t, actual.ID)
	require.Equal(t, *expected.ID, *actual.ID)

	// email
	require.NotNil(t, expected.Title)
	require.NotNil(t, actual.Title)
	require.Equal(t, *expected.Title, *actual.Title)

	// displayName
	require.NotNil(t, expected.Contents)
	require.NotNil(t, actual.Contents)
	require.Equal(t, *expected.Contents, *actual.Contents)

	// creation
	require.NotNil(t, expected.Creation)
	require.NotNil(t, actual.Creation)
	require.Equal(t, *expected.Creation, *actual.Creation)
}
