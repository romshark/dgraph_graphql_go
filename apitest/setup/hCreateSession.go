package setup

import (
	"time"

	"github.com/romshark/dgraph_graphql_go/api/graph/gqlmod"
	"github.com/romshark/dgraph_graphql_go/store/errors"
	"github.com/stretchr/testify/require"
)

func (h Helper) createSession(
	expectedErrorCode errors.Code,
	email string,
	password string,
) *gqlmod.Session {
	t := h.c.t

	var result struct {
		CreateSession *gqlmod.Session `json:"createSession"`
	}
	checkErr(t, expectedErrorCode, h.c.QueryVar(
		`mutation (
			$email: String!
			$password: String!
		) {
			createSession(
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

	require.NotNil(t, result.CreateSession)
	require.True(t, len(*result.CreateSession.Key) > 1)
	require.NotNil(t, result.CreateSession.User)
	require.Equal(t, email, *result.CreateSession.User.Email)
	require.WithinDuration(
		t,
		time.Now(),
		*result.CreateSession.Creation,
		h.creationTimeTollerance,
	)

	return result.CreateSession
}

// CreateSession helps creating a new session and assumes success
func (ok AssumeSuccess) CreateSession(
	email string,
	password string,
) *gqlmod.Session {
	return ok.h.createSession("", email, password)
}

// CreateSession helps creating a new session
func (notOkay AssumeFailure) CreateSession(
	expectedErrorCode errors.Code,
	email string,
	password string,
) {
	notOkay.checkErrCode(expectedErrorCode)
	notOkay.h.createSession(expectedErrorCode, email, password)
}
