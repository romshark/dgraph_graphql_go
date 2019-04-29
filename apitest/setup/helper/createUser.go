package helper

import (
	"time"

	"github.com/romshark/dgraph_graphql_go/api"

	"github.com/romshark/dgraph_graphql_go/api/graph/gqlmod"
	"github.com/stretchr/testify/require"
)

func (h Helper) createUser(
	successAssumption successAssumption,
	displayName string,
	email string,
) (*gqlmod.User, *api.ResponseError) {
	t := h.ts.T()

	var result struct {
		CreateUser *gqlmod.User `json:"createUser"`
	}
	err := h.ts.QueryVar(
		`mutation (
			$email: String!
			$displayName: String!
		) {
			createUser(
				email: $email
				displayName: $displayName,
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
		map[string]string{
			"displayName": displayName,
			"email":       email,
		},
		&result,
	)

	if successAssumption {
		require.Nil(t, err, 0)
	} else if err != nil {
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
	displayName string,
	email string,
) (*gqlmod.User, *api.ResponseError) {
	return h.createUser(potentialFailure, displayName, email)
}

// CreateUser helps creating a user and assumes success
func (ok AssumeSuccess) CreateUser(
	displayName string,
	email string,
) *gqlmod.User {
	result, _ := ok.h.createUser(success, displayName, email)
	return result
}
