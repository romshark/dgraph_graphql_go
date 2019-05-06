package setup

import (
	"time"

	"github.com/romshark/dgraph_graphql_go/api/graph"
	"github.com/romshark/dgraph_graphql_go/api/graph/gqlmod"
	"github.com/stretchr/testify/require"
)

func (h Helper) signIn(
	assumedSuccess successAssumption,
	email string,
	password string,
) (*gqlmod.Session, *graph.ResponseError) {
	t := h.c.t

	var result struct {
		SignIn *gqlmod.Session `json:"signIn"`
	}
	err := h.c.QueryVar(
		`mutation (
			$email: String!
			$password: String!
		) {
			signIn(
				email: $email
				password: $password,
			) {
				key
				creation
				user {
					id
					email
					displayName
					creation
				}
			}
		}`,
		map[string]interface{}{
			"email":    email,
			"password": password,
		},
		&result,
	)

	if err := checkErr(t, assumedSuccess, err); err != nil {
		return nil, err
	}

	require.NotNil(t, result.SignIn)
	require.True(t, len(*result.SignIn.Key) > 1)
	require.NotNil(t, result.SignIn.User)
	require.Equal(t, email, *result.SignIn.User.Email)
	require.WithinDuration(
		t,
		time.Now(),
		*result.SignIn.Creation,
		h.creationTimeTollerance,
	)

	return result.SignIn, nil
}

// SignIn helps signing in
func (h Helper) SignIn(
	email string,
	password string,
) (*gqlmod.Session, *graph.ResponseError) {
	return h.signIn(potentialFailure, email, password)
}

// SignIn helps signing in and assumes success
func (ok AssumeSuccess) SignIn(
	email string,
	password string,
) *gqlmod.Session {
	result, _ := ok.h.signIn(success, email, password)
	return result
}
