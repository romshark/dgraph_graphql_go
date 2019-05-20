package setup

import (
	"github.com/romshark/dgraph_graphql_go/api/graph/gqlmod"
	"github.com/romshark/dgraph_graphql_go/store/errors"
	"github.com/stretchr/testify/require"
)

func (h Helper) authenticate(
	expectedErrorCode errors.Code,
	sessionKey string,
) *gqlmod.Session {
	t := h.c.t

	var result struct {
		Authenticate *gqlmod.Session `json:"authenticate"`
	}
	checkErr(t, expectedErrorCode, h.c.QueryVar(
		`mutation (
			$sessionKey: String!
		) {
			authenticate(
				sessionKey: $sessionKey
			)
		}`,
		map[string]interface{}{
			"sessionKey": sessionKey,
		},
		&result,
	))

	if expectedErrorCode != "" {
		return nil
	}

	require.NotNil(t, result.Authenticate)

	return result.Authenticate
}

// Authenticate helps closing a session and assumes success
func (ok AssumeSuccess) Authenticate(
	sessionKey string,
) *gqlmod.Session {
	return ok.h.authenticate("", sessionKey)
}

// Authenticate helps closing a session
func (notOk AssumeFailure) Authenticate(
	expectedErrorCode errors.Code,
	sessionKey string,
) {
	notOk.checkErrCode(expectedErrorCode)
	notOk.h.authenticate(expectedErrorCode, sessionKey)
}
