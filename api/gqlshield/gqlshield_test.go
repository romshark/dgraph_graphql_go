package gqlshield_test

import (
	"testing"

	"github.com/romshark/dgraph_graphql_go/api/gqlshield"
	"github.com/stretchr/testify/require"
)

// TestWhitelisting tests WhitelistQueries
func TestWhitelisting(t *testing.T) {
	// Create a new shield instance
	shield, err := gqlshield.NewGraphQLShield(
		gqlshield.Config{},
		gqlshield.ClientRole{ID: 0, Name: "first"},
		gqlshield.ClientRole{ID: 1, Name: "second"},
	)
	require.NoError(t, err)
	require.NotNil(t, shield)

	// Define whitelist
	query1, err := shield.WhitelistQueries(gqlshield.Entry{
		Query: `query {
			users {
				id
				displayName
				email
			}
		}`,
		Name: "query one",
		Parameters: map[string]gqlshield.Parameter{
			"var1": gqlshield.Parameter{MaxValueLength: 1024},
		},
		WhitelistedFor: []int{0},
	})
	require.NoError(t, err)
	require.Len(t, query1, 1)
	require.NotNil(t, query1[0])

	query2, err := shield.WhitelistQueries(gqlshield.Entry{
		Query: `query {
			posts {
				id
				title
			}
		}`,
		Name:           "query two",
		WhitelistedFor: []int{0, 1},
	})
	require.NoError(t, err)
	require.Len(t, query2, 1)
	require.NotNil(t, query2[0])

	// Check
	err = shield.Check(
		0,
		query1[0].Query(),
		map[string]string{"var1": "v"},
	)
	require.NoError(t, err)

	err = shield.Check(
		0,
		query2[0].Query(),
		nil,
	)
	require.NoError(t, err)

	queries, err := shield.ListQueries()
	require.NoError(t, err)
	require.Len(t, queries, 2)
}

// TestWhitelistingErr tests all possible shield.WhitelistQueries errors
func TestWhitelistingErr(t *testing.T) {
	setup := func(t *testing.T) (shield gqlshield.GraphQLShield) {
		// Create a new shield instance
		var err error
		shield, err = gqlshield.NewGraphQLShield(
			gqlshield.Config{},
			gqlshield.ClientRole{ID: 0, Name: "first"},
			gqlshield.ClientRole{ID: 1, Name: "second"},
		)
		require.NoError(t, err)
		require.NotNil(t, shield)
		return
	}

	t.Run("duplicateQuery", func(t *testing.T) {
		shield := setup(t)

		query1, err := shield.WhitelistQueries(gqlshield.Entry{
			Query: `query {
				users {
					id
					displayName
					email
				}
			}`,
			Name: "query one",
			Parameters: map[string]gqlshield.Parameter{
				"var1": gqlshield.Parameter{MaxValueLength: 1024},
			},
			WhitelistedFor: []int{0},
		})
		require.NoError(t, err)
		require.NotNil(t, query1)

		query2, err := shield.WhitelistQueries(gqlshield.Entry{
			Query: `query {
				users {
					id
					displayName
					email
				}
			}`, // Duplicate query
			Name: "query two",
			Parameters: map[string]gqlshield.Parameter{
				"var1": gqlshield.Parameter{MaxValueLength: 1024},
			},
			WhitelistedFor: []int{0},
		})
		require.Error(t, err)
		require.Nil(t, query2)
	})

	t.Run("duplicateQueryName", func(t *testing.T) {
		shield := setup(t)

		query1, err := shield.WhitelistQueries(gqlshield.Entry{
			Query: `query { users { id displayName email } }`,
			Name:  "query one",
			Parameters: map[string]gqlshield.Parameter{
				"var1": gqlshield.Parameter{MaxValueLength: 1024},
			},
			WhitelistedFor: []int{0},
		})
		require.NoError(t, err)
		require.NotNil(t, query1)

		query2, err := shield.WhitelistQueries(gqlshield.Entry{
			Query: `query { users { id } }`,
			Name:  "query one", // Duplicate name
			Parameters: map[string]gqlshield.Parameter{
				"var1": gqlshield.Parameter{MaxValueLength: 1024},
			},
			WhitelistedFor: []int{0},
		})
		require.Error(t, err)
		require.Nil(t, query2)
	})

	t.Run("invalidRoles(no roles)", func(t *testing.T) {
		shield := setup(t)

		query, err := shield.WhitelistQueries(gqlshield.Entry{
			Query:          `query { users { id } }`,
			Name:           "query one",
			WhitelistedFor: []int{}, // Invalid: no roles
		})
		require.Error(t, err)
		require.Nil(t, query)
	})

	t.Run("invalidRoles(undefined roles)", func(t *testing.T) {
		shield := setup(t)

		query, err := shield.WhitelistQueries(gqlshield.Entry{
			Query:          `query { users { id } }`,
			Name:           "query one",
			WhitelistedFor: []int{0, 1, 999}, // Invalid: undefined role "999"
		})
		require.Error(t, err)
		require.Nil(t, query)
	})

	t.Run("invalidParameter(invalid name(empty))", func(t *testing.T) {
		shield := setup(t)

		query, err := shield.WhitelistQueries(gqlshield.Entry{
			Query: `query { users { id } }`,
			Name:  "query one",
			Parameters: map[string]gqlshield.Parameter{
				"": gqlshield.Parameter{ // Invalid parameter name
					MaxValueLength: 1024,
				},
			},
			WhitelistedFor: []int{0},
		})
		require.Error(t, err)
		require.Nil(t, query)
	})

	t.Run("invalidParameter(invalid MaxValueLength)", func(t *testing.T) {
		shield := setup(t)

		query, err := shield.WhitelistQueries(gqlshield.Entry{
			Query: `query { users { id } }`,
			Name:  "query one",
			Parameters: map[string]gqlshield.Parameter{
				"var1": gqlshield.Parameter{
					MaxValueLength: 0, // Invalid
				},
			},
			WhitelistedFor: []int{0},
		})
		require.Error(t, err)
		require.Nil(t, query)
	})
}

// TestRoleErr tests shield.WhitelistQueries
func TestRoleErr(t *testing.T) {
	// Create a new shield instance
	shield, err := gqlshield.NewGraphQLShield(
		gqlshield.Config{},
		gqlshield.ClientRole{ID: 1, Name: "first"},
		gqlshield.ClientRole{ID: 2, Name: "second"},
		gqlshield.ClientRole{ID: 3, Name: "third"},
		gqlshield.ClientRole{ID: 4, Name: "fourth"},
	)
	require.NoError(t, err)
	require.NotNil(t, shield)

	// Define whitelist
	queries, err := shield.WhitelistQueries(
		gqlshield.Entry{
			Query: `query {
				users {
					id
					displayName
					email
				}
			}`,
			Name: "query one",
			Parameters: map[string]gqlshield.Parameter{
				"var1": gqlshield.Parameter{MaxValueLength: 1024},
			},
			WhitelistedFor: []int{1},
		},
		gqlshield.Entry{
			Query: `query {
				posts {
					id
					title
				}
			}`,
			Name:           "query two",
			WhitelistedFor: []int{1, 2},
		},
		gqlshield.Entry{
			Query: `query( $userID: Identifier! ) {
				user( userID: $userID ) {
					id
					displayName
					posts {
						id
						title
					}
				}
			}`,
			Name: "query three",
			Parameters: map[string]gqlshield.Parameter{
				"userID":        gqlshield.Parameter{MaxValueLength: 32},
				"postListLimit": gqlshield.Parameter{MaxValueLength: 8},
			},
			WhitelistedFor: []int{4},
		},
	)
	require.NoError(t, err)
	require.Len(t, queries, 3)
	require.NotNil(t, queries[0])
	require.NotNil(t, queries[1])
	require.NotNil(t, queries[2])

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

	check(queries[0], map[string]string{"var1": "v"}, Expect{
		1: true,
		2: false,
		3: false,
		4: false,
	})
	check(queries[1], nil, Expect{
		1: true,
		2: true,
		3: false,
		4: false,
	})
	check(queries[2], map[string]string{
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
		gqlshield.Config{},
		gqlshield.ClientRole{ID: 0, Name: "default"},
	)
	require.NoError(t, err)
	require.NotNil(t, shield)

	// Define whitelist
	queries, err := shield.WhitelistQueries(
		gqlshield.Entry{
			Query: `query {
				users {
					id
					displayName
					email
				}
			}`,
			Name: "query one",
			Parameters: map[string]gqlshield.Parameter{
				"var1": gqlshield.Parameter{MaxValueLength: 1024},
			},
			WhitelistedFor: []int{0},
		},
		gqlshield.Entry{
			Query: `query {
				posts {
					id
					title
				}
			}`,
			Name:           "query two",
			WhitelistedFor: []int{0},
		},
	)
	require.NoError(t, err)
	require.Len(t, queries, 2)
	require.NotNil(t, queries[0])
	require.NotNil(t, queries[1])

	// Remove first query
	require.NoError(t, shield.RemoveQuery(queries[0]))

	// Check
	err = shield.Check(
		0,
		queries[0].Query(),
		map[string]string{"var1": "v"},
	)
	require.Error(t, err)
	require.Equal(t, gqlshield.ErrUnauthorized, gqlshield.ErrCode(err))

	err = shield.Check(
		0,
		queries[1].Query(),
		nil,
	)
	require.NoError(t, err)

	listedQueries, err := shield.ListQueries()
	require.NoError(t, err)
	require.Len(t, listedQueries, 1)
	require.Equal(t, queries[1].Name(), listedQueries["query two"].Name())

	// Remove second query
	require.NoError(t, shield.RemoveQuery(queries[1]))

	// Check
	err = shield.Check(
		0,
		queries[1].Query(),
		nil,
	)
	require.Error(t, err)
	require.Equal(t, gqlshield.ErrUnauthorized, gqlshield.ErrCode(err))

	listedQueries, err = shield.ListQueries()
	require.NoError(t, err)
	require.Len(t, listedQueries, 0)
}

// TestWrongArg tests argument validation
func TestWrongArg(t *testing.T) {
	setup := func() (shield gqlshield.GraphQLShield, query gqlshield.Query) {
		// Create a new shield instance
		var err error
		shield, err = gqlshield.NewGraphQLShield(
			gqlshield.Config{},
			gqlshield.ClientRole{ID: 0, Name: "default"},
		)
		require.NoError(t, err)
		require.NotNil(t, shield)

		queries, err := shield.WhitelistQueries(gqlshield.Entry{
			Query: `query($id: String!) {
				users( id: $id ) {
					name
				}
			}`,
			Name: "user name",
			Parameters: map[string]gqlshield.Parameter{
				"id": gqlshield.Parameter{MaxValueLength: 32},
			},
			WhitelistedFor: []int{0},
		})
		require.NoError(t, err)
		require.Len(t, queries, 1)

		query = queries[0]
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
		gqlshield.Config{},
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
		shield, err := gqlshield.NewGraphQLShield(gqlshield.Config{})
		require.Error(t, err)
		require.Nil(t, shield)
	})

	t.Run("duplicateRoleID", func(t *testing.T) {
		shield, err := gqlshield.NewGraphQLShield(
			gqlshield.Config{},
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
			gqlshield.Config{},
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
			gqlshield.Config{},
			gqlshield.ClientRole{
				ID:   1,
				Name: "",
			},
		)
		require.Error(t, err)
		require.Nil(t, shield)
	})
}

// TestQueries tests shield.ListQueries
func TestQueries(t *testing.T) {
	// Create and initialize a new shield instance
	shield, err := gqlshield.NewGraphQLShield(
		gqlshield.Config{},
		gqlshield.ClientRole{ID: 0, Name: "first"},
		gqlshield.ClientRole{ID: 1, Name: "second"},
	)
	require.NoError(t, err)
	require.NotNil(t, shield)

	// Define whitelist
	queries, err := shield.WhitelistQueries(
		gqlshield.Entry{
			Query: `query { users { id displayName email } }`,
			Name:  "query one",
			Parameters: map[string]gqlshield.Parameter{
				"var1": gqlshield.Parameter{MaxValueLength: 1024},
			},
			WhitelistedFor: []int{0},
		},
		gqlshield.Entry{
			Query:          `query { posts { id title } }`,
			Name:           "query two",
			WhitelistedFor: []int{0, 1},
		},
	)
	require.NoError(t, err)
	require.Len(t, queries, 2)

	// Check
	listedQueries, err := shield.ListQueries()
	require.NoError(t, err)
	require.Len(t, listedQueries, 2)

	q1, exists := listedQueries["query one"]
	require.True(t, exists)
	require.NotNil(t, q1)
	require.Equal(t, queries[0].Name(), q1.Name())
	require.Equal(t, queries[0].Query(), q1.Query())
	require.Equal(t, queries[0].Parameters(), q1.Parameters())
	require.Equal(t, queries[0].WhitelistedFor(), q1.WhitelistedFor())

	q2, exists := listedQueries["query two"]
	require.True(t, exists)
	require.NotNil(t, q1)
	require.Equal(t, queries[1].Name(), q2.Name())
	require.Equal(t, queries[1].Query(), q2.Query())
	require.Equal(t, queries[1].Parameters(), q2.Parameters())
	require.Equal(t, queries[1].WhitelistedFor(), q2.WhitelistedFor())

	// Remove first query
	require.NoError(t, shield.RemoveQuery(queries[0]))
	listedQueries, err = shield.ListQueries()
	require.NoError(t, err)

	q1, exists = listedQueries["query one"]
	require.False(t, exists)
	require.Nil(t, q1)

	q2, exists = listedQueries["query two"]
	require.True(t, exists)
	require.NotNil(t, q2)
	require.Equal(t, queries[1].Name(), q2.Name())
	require.Equal(t, queries[1].Query(), q2.Query())
	require.Equal(t, queries[1].Parameters(), q2.Parameters())
	require.Equal(t, queries[1].WhitelistedFor(), q2.WhitelistedFor())

	// Remove second query
	require.NoError(t, shield.RemoveQuery(queries[1]))
	listedQueries, err = shield.ListQueries()
	require.NoError(t, err)

	q1, exists = listedQueries["query one"]
	require.False(t, exists)
	require.Nil(t, q1)

	q2, exists = listedQueries["query two"]
	require.False(t, exists)
	require.Nil(t, q2)
}

type persistencyManagerMock struct {
	loads int
	saves []*gqlshield.State
}

func (m *persistencyManagerMock) lastSave() *gqlshield.State {
	if len(m.saves) < 1 {
		return nil
	}
	return m.saves[len(m.saves)-1]
}

func (m *persistencyManagerMock) Load() (*gqlshield.State, error) {
	m.loads++
	return nil, nil
}

func (m *persistencyManagerMock) Save(state *gqlshield.State) error {
	m.saves = append(m.saves, state)
	return nil
}

// TestPersistency tests the persistency manager option
func TestPersistency(t *testing.T) {
	// Create a new persistency manager mock instance
	mock := &persistencyManagerMock{
		loads: 0,
		saves: make([]*gqlshield.State, 0),
	}

	// Create a new shield instance
	shield, err := gqlshield.NewGraphQLShield(
		gqlshield.Config{
			PersistencyManager: mock,
		},
		gqlshield.ClientRole{ID: 1, Name: "first"},
	)
	require.NoError(t, err)
	require.NotNil(t, shield)

	// Define helper function
	check := func(
		expectedLoads int,
		expectedSaves int,
		expectedEntries ...gqlshield.Entry,
	) {
		require.Equal(t, mock.loads, expectedLoads)
		require.Len(t, mock.saves, expectedSaves)

		// Compare the expected and actual states
		// if any saves were expected to have happened
		if expectedSaves > 0 {
			actualState := mock.lastSave()

			// Compare the client roles
			require.Equal(t, []gqlshield.ClientRole{
				gqlshield.ClientRole{
					ID: 1, Name: "first",
				},
			}, actualState.Roles)
			require.Len(
				t,
				actualState.WhitelistedQueries,
				len(expectedEntries),
			)

			// Compare expected queries with actually saved queries
			findQueryByName := func(name string) *gqlshield.QueryModel {
				for _, query := range actualState.WhitelistedQueries {
					if query.Name == name {
						return &query
					}
				}
				return nil
			}
			for _, expected := range expectedEntries {
				// Find by name
				actual := findQueryByName(expected.Name)
				require.NotNil(t, actual)

				require.Equal(t, expected.Name, actual.Name)
				require.Equal(t, expected.Query, actual.Query)
				require.Equal(t, expected.Parameters, actual.Parameters)
				require.Equal(t, expected.WhitelistedFor, actual.WhitelistedFor)
			}
		} else {
			require.Nil(t, mock.lastSave())
		}
	}

	// Expect the persistence manager to have loaded the initial state
	check(1, 0)

	// Add the first query
	query1Entry := gqlshield.Entry{
		Query: `query( $uid: ID! ) { user( id: $uid ) { name email } }`,
		Name:  "query one",
		Parameters: map[string]gqlshield.Parameter{
			"var1": gqlshield.Parameter{MaxValueLength: 1024},
		},
		WhitelistedFor: []int{1},
	}
	query1, err := shield.WhitelistQueries(query1Entry)
	require.NoError(t, err)
	require.NotNil(t, query1)

	// Expect the persistence manager to have saved the state including query 1
	check(1, 1, query1Entry)

	// Add a second query
	query2Entry := gqlshield.Entry{
		Query:          `query { posts { id title } }`,
		Name:           "query two",
		WhitelistedFor: []int{1},
	}
	query2, err := shield.WhitelistQueries(query2Entry)
	require.NoError(t, err)
	require.NotNil(t, query2)

	// Expect the persistence manager to have saved the state again
	// including query 2
	check(1, 2, query1Entry, query2Entry)

	// Remove first query
	require.NoError(t, shield.RemoveQuery(query1[0]))

	// Expect the persistence manager to have saved the state again
	// not including query 1 any longer
	check(1, 3, query2Entry)

	// Remove second query
	require.NoError(t, shield.RemoveQuery(query2[0]))

	// Expect the persistence manager to have saved the state again
	// not including any query
	check(1, 4)
}
