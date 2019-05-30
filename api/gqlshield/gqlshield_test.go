package gqlshield_test

import (
	"testing"

	"github.com/romshark/dgraph_graphql_go/api/gqlshield"
	"github.com/stretchr/testify/require"
)

func TestGQLShield(t *testing.T) {
	query1 := gqlshield.Query{
		Query: []byte(`query {
			users {
				id
				displayName
				email
			}
		}`),
		Name: "query one",
		Parameters: map[string]gqlshield.Parameter{
			"var1": gqlshield.Parameter{MaxValueLength: 1024},
		},
	}
	query2 := gqlshield.Query{
		Query: []byte(`query {
			posts {
				id
				title
			}
		}`),
		Name: "query two",
	}

	// Create a new shield instance
	shield := gqlshield.NewGraphQLShield()
	require.NotNil(t, shield)

	q := query1.Clone()
	query1isIn, err := shield.Check(q.Query, map[string]string{"var1": "v"})
	require.NoError(t, err)
	require.False(t, query1isIn)

	q = query2.Clone()
	query2isIn, err := shield.Check(q.Query, nil)
	require.NoError(t, err)
	require.False(t, query2isIn)

	require.Len(t, shield.Queries(), 0)

	// Insert first query
	require.NoError(t, shield.WhitelistQuery(query1.Clone()))

	q = query1.Clone()
	query1isIn2, err := shield.Check(q.Query, map[string]string{"var1": "v"})
	require.NoError(t, err)
	require.True(t, query1isIn2)

	q = query2.Clone()
	query2isIn2, err := shield.Check(q.Query, nil)
	require.NoError(t, err)
	require.False(t, query2isIn2)

	require.Len(t, shield.Queries(), 1)

	// Insert second query
	require.NoError(t, shield.WhitelistQuery(query2.Clone()))

	q = query1.Clone()
	query1isIn3, err := shield.Check(q.Query, map[string]string{"var1": "v"})
	require.NoError(t, err)
	require.True(t, query1isIn3)

	q = query2.Clone()
	query2isIn3, err := shield.Check(q.Query, nil)
	require.NoError(t, err)
	require.True(t, query2isIn3)

	require.Len(t, shield.Queries(), 2)

	// Remove first query
	q = query1.Clone()
	removed, err := shield.RemoveQuery(q.Query)
	require.NoError(t, err)
	require.NotNil(t, removed)

	q = query1.Clone()
	query1isIn4, err := shield.Check(q.Query, map[string]string{"var1": "v"})
	require.NoError(t, err)
	require.False(t, query1isIn4)

	q = query2.Clone()
	query2isIn4, err := shield.Check(q.Query, nil)
	require.NoError(t, err)
	require.True(t, query2isIn4)

	require.Len(t, shield.Queries(), 1)

	// Remove second query
	q = query2.Clone()
	removed, err = shield.RemoveQuery(q.Query)
	require.NoError(t, err)
	require.NotNil(t, removed)

	q = query1.Clone()
	query1isIn5, err := shield.Check(q.Query, map[string]string{"var1": "v"})
	require.NoError(t, err)
	require.False(t, query1isIn5)

	q = query2.Clone()
	query2isIn5, err := shield.Check(q.Query, nil)
	require.NoError(t, err)
	require.False(t, query2isIn5)

	require.Len(t, shield.Queries(), 0)
}

func TestGQLShieldWrongArg(t *testing.T) {
	setup := func() (shield gqlshield.GraphQLShield, query *gqlshield.Query) {
		// Create a new shield instance
		shield = gqlshield.NewGraphQLShield()
		require.NotNil(t, shield)

		query = &gqlshield.Query{
			Query: []byte(`query($id: String!) {
				users( id: $id ) {
					name
				}
			}`),
			Name: "user name",
			Parameters: map[string]gqlshield.Parameter{
				"id": gqlshield.Parameter{MaxValueLength: 32},
			},
		}

		require.NoError(t, shield.WhitelistQuery(query.Clone()))
		return
	}

	t.Run("wrongName", func(t *testing.T) {
		shield, qr := setup()
		found, err := shield.Check(
			qr.Query,
			map[string]string{"wrongName": "v"},
		)
		require.True(t, found)
		require.Error(t, err)
	})

	t.Run("maxLenExceeded", func(t *testing.T) {
		shield, qr := setup()
		found, err := shield.Check(
			qr.Query,
			map[string]string{"wrongName": "11110000111100001111000011110000F"},
		)
		require.True(t, found)
		require.Error(t, err)
	})
}
