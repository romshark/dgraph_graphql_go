package setup

import (
	"time"

	"github.com/romshark/dgraph_graphql_go/api/graph/gqlmod"
	"github.com/romshark/dgraph_graphql_go/store/errors"
	"github.com/stretchr/testify/require"
)

func (h Helper) createUser(
	expectedErrorCode errors.Code,
	displayName,
	email,
	password string,
) *gqlmod.User {
	t := h.c.t

	var result struct {
		CreateUser *gqlmod.User `json:"createUser"`
	}
	checkErr(t, expectedErrorCode, h.c.QueryVar(
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
	))

	if expectedErrorCode != "" {
		return nil
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

	return result.CreateUser
}

// CreateUser helps creating a user and assumes success
func (ok AssumeSuccess) CreateUser(
	displayName,
	email,
	password string,
) *gqlmod.User {
	return ok.h.createUser("", displayName, email, password)
}

// CreateUser helps creating a user
func (notOkay AssumeFailure) CreateUser(
	expectedErrorCode errors.Code,
	displayName,
	email,
	password string,
) {
	notOkay.checkErrCode(expectedErrorCode)
	notOkay.h.createUser(expectedErrorCode, displayName, email, password)
}
