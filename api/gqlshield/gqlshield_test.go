package gqlshield_test

import (
	"testing"

	"github.com/romshark/dgraph_graphql_go/api/gqlshield"
	"github.com/stretchr/testify/require"
)

// TestWhitelisting tests WhitelistQuery
func TestWhitelisting(t *testing.T) {
	// Create a new shield instance
	shield, err := gqlshield.NewGraphQLShield(
		gqlshield.ClientRole{ID: 0, Name: "first"},
		gqlshield.ClientRole{ID: 1, Name: "second"},
	)
	require.NoError(t, err)
	require.NotNil(t, shield)

	// Define whitelist
	query1, err := shield.WhitelistQuery(
		[]byte(`query {
			users {
				id
				displayName
				email
			}
		}`),
		"query one",
		map[string]gqlshield.Parameter{
			"var1": gqlshield.Parameter{MaxValueLength: 1024},
		},
		[]int{0},
	)
	require.NoError(t, err)
	require.NotNil(t, query1)

	query2, err := shield.WhitelistQuery(
		[]byte(`query {
			posts {
				id
				title
			}
		}`),
		"query two",
		nil,
		[]int{0, 1},
	)
	require.NoError(t, err)
	require.NotNil(t, query2)

	// Check
	err = shield.Check(
		0,
		query1.Query(),
		map[string]string{"var1": "v"},
	)
	require.NoError(t, err)

	err = shield.Check(
		0,
		query2.Query(),
		nil,
	)
	require.NoError(t, err)

	queries, err := shield.Queries()
	require.NoError(t, err)
	require.Len(t, queries, 2)
}

// TestWhitelistingErr tests all possible shield.WhitelistQuery errors
func TestWhitelistingErr(t *testing.T) {
	setup := func(t *testing.T) (shield gqlshield.GraphQLShield) {
		// Create a new shield instance
		var err error
		shield, err = gqlshield.NewGraphQLShield(
			gqlshield.ClientRole{ID: 0, Name: "first"},
			gqlshield.ClientRole{ID: 1, Name: "second"},
		)
		require.NoError(t, err)
		require.NotNil(t, shield)
		return
	}

	t.Run("duplicateQuery", func(t *testing.T) {
		shield := setup(t)

		query1, err := shield.WhitelistQuery(
			[]byte(`query {
				users {
					id
					displayName
					email
				}
			}`),
			"query one",
			map[string]gqlshield.Parameter{
				"var1": gqlshield.Parameter{MaxValueLength: 1024},
			},
			[]int{0},
		)
		require.NoError(t, err)
		require.NotNil(t, query1)

		query2, err := shield.WhitelistQuery(
			[]byte(`query {
				users {
					id
					displayName
					email
				}
			}`), // Duplicate query
			"query two",
			map[string]gqlshield.Parameter{
				"var1": gqlshield.Parameter{MaxValueLength: 1024},
			},
			[]int{0},
		)
		require.Error(t, err)
		require.Nil(t, query2)
	})

	t.Run("duplicateQueryName", func(t *testing.T) {
		shield := setup(t)

		query1, err := shield.WhitelistQuery(
			[]byte(`query { users { id displayName email } }`),
			"query one",
			map[string]gqlshield.Parameter{
				"var1": gqlshield.Parameter{MaxValueLength: 1024},
			},
			[]int{0},
		)
		require.NoError(t, err)
		require.NotNil(t, query1)

		query2, err := shield.WhitelistQuery(
			[]byte(`query { users { id } }`),
			"query one", // Duplicate name
			map[string]gqlshield.Parameter{
				"var1": gqlshield.Parameter{MaxValueLength: 1024},
			},
			[]int{0},
		)
		require.Error(t, err)
		require.Nil(t, query2)
	})

	t.Run("invalidRoles(no roles)", func(t *testing.T) {
		shield := setup(t)

		query, err := shield.WhitelistQuery(
			[]byte(`query { users { id } }`),
			"query one",
			nil,
			[]int{}, // Invalid: no roles
		)
		require.Error(t, err)
		require.Nil(t, query)
	})

	t.Run("invalidRoles(undefined roles)", func(t *testing.T) {
		shield := setup(t)

		query, err := shield.WhitelistQuery(
			[]byte(`query { users { id } }`),
			"query one",
			nil,
			[]int{0, 1, 999}, // Invalid: undefined role "999"
		)
		require.Error(t, err)
		require.Nil(t, query)
	})

	t.Run("invalidParameter(invalid name(empty))", func(t *testing.T) {
		shield := setup(t)

		query, err := shield.WhitelistQuery(
			[]byte(`query { users { id } }`),
			"query one",
			map[string]gqlshield.Parameter{
				"": gqlshield.Parameter{ // Invalid parameter name
					MaxValueLength: 1024,
				},
			},
			[]int{0},
		)
		require.Error(t, err)
		require.Nil(t, query)
	})

	t.Run("invalidParameter(invalid MaxValueLength)", func(t *testing.T) {
		shield := setup(t)

		query, err := shield.WhitelistQuery(
			[]byte(`query { users { id } }`),
			"query one",
			map[string]gqlshield.Parameter{
				"var1": gqlshield.Parameter{
					MaxValueLength: 0, // Invalid
				},
			},
			[]int{0},
		)
		require.Error(t, err)
		require.Nil(t, query)
	})
}

// TestRoleErr tests shield.WhitelistQuery
func TestRoleErr(t *testing.T) {
	// Create a new shield instance
	shield, err := gqlshield.NewGraphQLShield(
		gqlshield.ClientRole{ID: 1, Name: "first"},
		gqlshield.ClientRole{ID: 2, Name: "second"},
		gqlshield.ClientRole{ID: 3, Name: "third"},
		gqlshield.ClientRole{ID: 4, Name: "fourth"},
	)
	require.NoError(t, err)
	require.NotNil(t, shield)

	// Define whitelist
	query1, err := shield.WhitelistQuery(
		[]byte(`query {
			users {
				id
				displayName
				email
			}
		}`),
		"query one",
		map[string]gqlshield.Parameter{
			"var1": gqlshield.Parameter{MaxValueLength: 1024},
		},
		[]int{1},
	)
	require.NoError(t, err)
	require.NotNil(t, query1)

	query2, err := shield.WhitelistQuery(
		[]byte(`query {
			posts {
				id
				title
			}
		}`),
		"query two",
		nil,
		[]int{1, 2},
	)
	require.NoError(t, err)
	require.NotNil(t, query2)

	query3, err := shield.WhitelistQuery(
		[]byte(`query( $userID: Identifier! ) {
			user( userID: $userID ) {
				id
				displayName
				posts {
					id
					title
				}
			}
		}`),
		"query three",
		map[string]gqlshield.Parameter{
			"userID":        gqlshield.Parameter{MaxValueLength: 32},
			"postListLimit": gqlshield.Parameter{MaxValueLength: 8},
		},
		[]int{4},
	)
	require.NoError(t, err)
	require.NotNil(t, query3)

	type Expect map[int]bool
	check := func(
		query gqlshield.Query,
		args map[string]string,
		expectancy Expect,
	) {
		for role, expectAuth := range expectancy {
			err := shield.Check(
				role,
				query.Query(),
				args,
			)
			if expectAuth {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				require.Equal(
					t,
					gqlshield.ErrUnauthorized,
					gqlshield.ErrCode(err),
				)
			}
		}
	}

	check(query1, map[string]string{"var1": "v"}, Expect{
		1: true,
		2: false,
		3: false,
		4: false,
	})
	check(query2, nil, Expect{
		1: true,
		2: true,
		3: false,
		4: false,
	})
	check(query3, map[string]string{
		"userID":        "12345678901234567890123456789012",
		"postListLimit": "50",
	}, Expect{
		1: false,
		2: false,
		3: false,
		4: true,
	})
}

// TestRemove tests shield.RemoveQuery
func TestRemove(t *testing.T) {
	// Create a new shield instance
	shield, err := gqlshield.NewGraphQLShield(
		gqlshield.ClientRole{ID: 0, Name: "default"},
	)
	require.NoError(t, err)
	require.NotNil(t, shield)

	// Define whitelist
	query1, err := shield.WhitelistQuery(
		[]byte(`query {
			users {
				id
				displayName
				email
			}
		}`),
		"query one",
		map[string]gqlshield.Parameter{
			"var1": gqlshield.Parameter{MaxValueLength: 1024},
		},
		[]int{0},
	)
	require.NoError(t, err)
	require.NotNil(t, query1)

	query2, err := shield.WhitelistQuery(
		[]byte(`query {
			posts {
				id
				title
			}
		}`),
		"query two",
		nil,
		[]int{0},
	)
	require.NoError(t, err)
	require.NotNil(t, query2)

	// Remove first query
	require.NoError(t, shield.RemoveQuery(query1))

	// Check
	err = shield.Check(
		0,
		query1.Query(),
		map[string]string{"var1": "v"},
	)
	require.Error(t, err)
	require.Equal(t, gqlshield.ErrUnauthorized, gqlshield.ErrCode(err))

	err = shield.Check(
		0,
		query2.Query(),
		nil,
	)
	require.NoError(t, err)

	queries, err := shield.Queries()
	require.NoError(t, err)
	require.Len(t, queries, 1)
	require.Equal(t, query2.Name(), queries["query two"].Name())

	// Remove second query
	require.NoError(t, shield.RemoveQuery(query2))

	// Check
	err = shield.Check(
		0,
		query2.Query(),
		nil,
	)
	require.Error(t, err)
	require.Equal(t, gqlshield.ErrUnauthorized, gqlshield.ErrCode(err))

	queries, err = shield.Queries()
	require.NoError(t, err)
	require.Len(t, queries, 0)
}

// TestWrongArg tests argument validation
func TestWrongArg(t *testing.T) {
	setup := func() (shield gqlshield.GraphQLShield, query gqlshield.Query) {
		// Create a new shield instance
		var err error
		shield, err = gqlshield.NewGraphQLShield(
			gqlshield.ClientRole{ID: 0, Name: "default"},
		)
		require.NoError(t, err)
		require.NotNil(t, shield)

		query, err = shield.WhitelistQuery(
			[]byte(`query($id: String!) {
				users( id: $id ) {
					name
				}
			}`),
			"user name",
			map[string]gqlshield.Parameter{
				"id": gqlshield.Parameter{MaxValueLength: 32},
			},
			[]int{0},
		)
		require.NoError(t, err)
		require.NotNil(t, query)
		return
	}

	t.Run("wrongName", func(t *testing.T) {
		shield, qr := setup()
		err := shield.Check(
			0,
			qr.Query(),
			map[string]string{"wrongName": "v"},
		)
		require.Error(t, err)
		require.Equal(t, gqlshield.ErrUnauthorized, gqlshield.ErrCode(err))
	})

	t.Run("maxLenExceeded", func(t *testing.T) {
		shield, qr := setup()
		err := shield.Check(
			0,
			qr.Query(),
			map[string]string{"wrongName": "11110000111100001111000011110000F"},
		)
		require.Error(t, err)
		require.Equal(t, gqlshield.ErrUnauthorized, gqlshield.ErrCode(err))
	})
}

// TestNewGraphQLShield tests the NewGraphQLShield constructor function
func TestNewGraphQLShield(t *testing.T) {
	shield, err := gqlshield.NewGraphQLShield(
		gqlshield.ClientRole{
			ID:   1,
			Name: "first",
		},
		gqlshield.ClientRole{
			ID:   2,
			Name: "second",
		},
		gqlshield.ClientRole{
			ID:   3,
			Name: "third",
		},
	)
	require.NoError(t, err)
	require.NotNil(t, shield)
}

// TestNewGraphQLShieldErr tests all possible errors for NewGraphQLShield
func TestNewGraphQLShieldErr(t *testing.T) {
	t.Run("noRoles", func(t *testing.T) {
		shield, err := gqlshield.NewGraphQLShield()
		require.Error(t, err)
		require.Nil(t, shield)
	})

	t.Run("duplicateRoleID", func(t *testing.T) {
		shield, err := gqlshield.NewGraphQLShield(
			gqlshield.ClientRole{
				ID:   1,
				Name: "first",
			},
			gqlshield.ClientRole{
				ID:   2,
				Name: "second",
			},
			gqlshield.ClientRole{
				ID:   1, // Duplicate client role ID
				Name: "third",
			},
		)
		require.Error(t, err)
		require.Nil(t, shield)
	})

	t.Run("duplicateRoleName", func(t *testing.T) {
		shield, err := gqlshield.NewGraphQLShield(
			gqlshield.ClientRole{
				ID:   1,
				Name: "first",
			},
			gqlshield.ClientRole{
				ID:   2,
				Name: "second",
			},
			gqlshield.ClientRole{
				ID:   3,
				Name: "first", // Duplicate client role name
			},
		)
		require.Error(t, err)
		require.Nil(t, shield)
	})

	t.Run("invalidRoleName(empty)", func(t *testing.T) {
		shield, err := gqlshield.NewGraphQLShield(
			gqlshield.ClientRole{
				ID:   1,
				Name: "",
			},
		)
		require.Error(t, err)
		require.Nil(t, shield)
	})
}

// TestQueries tests shield.Queries
func TestQueries(t *testing.T) {
	// Create a new shield instance
	shield, err := gqlshield.NewGraphQLShield(
		gqlshield.ClientRole{ID: 0, Name: "first"},
		gqlshield.ClientRole{ID: 1, Name: "second"},
	)
	require.NoError(t, err)
	require.NotNil(t, shield)

	// Define whitelist
	query1, err := shield.WhitelistQuery(
		[]byte(`query {
			users {
				id
				displayName
				email
			}
		}`),
		"query one",
		map[string]gqlshield.Parameter{
			"var1": gqlshield.Parameter{MaxValueLength: 1024},
		},
		[]int{0},
	)
	require.NoError(t, err)
	require.NotNil(t, query1)

	query2, err := shield.WhitelistQuery(
		[]byte(`query {
			posts {
				id
				title
			}
		}`),
		"query two",
		nil,
		[]int{0, 1},
	)
	require.NoError(t, err)
	require.NotNil(t, query2)

	// Check
	queries, err := shield.Queries()
	require.NoError(t, err)
	require.Len(t, queries, 2)

	q1, exists := queries["query one"]
	require.True(t, exists)
	require.NotNil(t, q1)
	require.Equal(t, query1.Name(), q1.Name())
	require.Equal(t, query1.Query(), q1.Query())
	require.Equal(t, query1.Parameters(), q1.Parameters())
	require.Equal(t, query1.WhitelistedFor(), q1.WhitelistedFor())

	q2, exists := queries["query two"]
	require.True(t, exists)
	require.NotNil(t, q1)
	require.Equal(t, query2.Name(), q2.Name())
	require.Equal(t, query2.Query(), q2.Query())
	require.Equal(t, query2.Parameters(), q2.Parameters())
	require.Equal(t, query2.WhitelistedFor(), q2.WhitelistedFor())

	// Remove first query
	require.NoError(t, shield.RemoveQuery(query1))
	queries, err = shield.Queries()
	require.NoError(t, err)

	q1, exists = queries["query one"]
	require.False(t, exists)
	require.Nil(t, q1)

	q2, exists = queries["query two"]
	require.True(t, exists)
	require.NotNil(t, q2)
	require.Equal(t, query2.Name(), q2.Name())
	require.Equal(t, query2.Query(), q2.Query())
	require.Equal(t, query2.Parameters(), q2.Parameters())
	require.Equal(t, query2.WhitelistedFor(), q2.WhitelistedFor())

	// Remove second query
	require.NoError(t, shield.RemoveQuery(query2))
	queries, err = shield.Queries()
	require.NoError(t, err)

	q1, exists = queries["query one"]
	require.False(t, exists)
	require.Nil(t, q1)

	q2, exists = queries["query two"]
	require.False(t, exists)
	require.Nil(t, q2)
}
