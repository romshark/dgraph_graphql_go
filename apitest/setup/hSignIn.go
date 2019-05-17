package setup

import (
	"time"

	"github.com/romshark/dgraph_graphql_go/api/graph/gqlmod"
	"github.com/romshark/dgraph_graphql_go/store/errors"
	"github.com/stretchr/testify/require"
)

func (h Helper) signIn(
	expectedErrorCode errors.Code,
	email string,
	password string,
) *gqlmod.Session {
	t := h.c.t

	var result struct {
		SignIn *gqlmod.Session `json:"signIn"`
	}
	checkErr(t, expectedErrorCode, h.c.QueryVar(
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
	))

	if expectedErrorCode != "" {
		return nil
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

	return result.SignIn
}

// SignIn helps signing in and assumes success
func (ok AssumeSuccess) SignIn(
	email string,
	password string,
) *gqlmod.Session {
	return ok.h.signIn("", email, password)
}

// SignIn helps signing in
func (notOkay AssumeFailure) SignIn(
	expectedErrorCode errors.Code,
	email string,
	password string,
) {
	notOkay.checkErrCode(expectedErrorCode)
	notOkay.h.signIn(expectedErrorCode, email, password)
}
