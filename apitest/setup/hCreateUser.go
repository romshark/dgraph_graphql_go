package setup

import (
	"time"

	"github.com/romshark/dgraph_graphql_go/api/graph"
	"github.com/romshark/dgraph_graphql_go/api/graph/gqlmod"
	"github.com/stretchr/testify/require"
)

func (h Helper) createUser(
	assumedSuccess successAssumption,
	displayName,
	email,
	password string,
) (*gqlmod.User, *graph.ResponseError) {
	t := h.c.t

	var result struct {
		CreateUser *gqlmod.User `json:"createUser"`
	}
	err := h.c.QueryVar(
		`mutation (
			$email: String!
			$displayName: String!
			$password: String!
		) {
			createUser(
				email: $email
				displayName: $displayName
				password: $password
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
		map[string]interface{}{
			"displayName": displayName,
			"email":       email,
			"password":    password,
		},
		&result,
	)

	if err := checkErr(t, assumedSuccess, err); err != nil {
		return nil, err
	}

	require.NotNil(t, result.CreateUser)
	require.Len(t, *result.CreateUser.ID, 32)
	require.Equal(t, email, *result.CreateUser.Email)
	require.Equal(t, displayName, *result.CreateUser.DisplayName)
	require.Len(t, result.CreateUser.Posts, 0)
	require.WithinDuration(
		t,
		time.Now(),
		*result.CreateUser.Creation,
		h.creationTimeTollerance,
	)

	return result.CreateUser, nil
}

// CreateUser helps creating a user
func (h Helper) CreateUser(
	displayName,
	email,
	password string,
) (*gqlmod.User, *graph.ResponseError) {
	return h.createUser(potentialFailure, displayName, email, password)
}

// CreateUser helps creating a user and assumes success
func (ok AssumeSuccess) CreateUser(
	displayName,
	email,
	password string,
) *gqlmod.User {
	result, _ := ok.h.createUser(success, displayName, email, password)
	return result
}
