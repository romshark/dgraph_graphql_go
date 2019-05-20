package setup

import (
	"github.com/romshark/dgraph_graphql_go/store"
	"github.com/romshark/dgraph_graphql_go/store/errors"
)

func (h Helper) closeAllSessions(
	expectedErrorCode errors.Code,
	user store.ID,
) []string {
	t := h.c.t

	var result struct {
		CloseAllSessions []string `json:"closeAllSessions"`
	}
	checkErr(t, expectedErrorCode, h.c.QueryVar(
		`mutation (
			$user: Identifier!
		) {
			closeAllSessions(
				user: $user
			)
		}`,
		map[string]interface{}{
			"user": string(user),
		},
		&result,
	))

	if expectedErrorCode != "" {
		return nil
	}

	return result.CloseAllSessions
}

// CloseAllSessions helps closing all sessions of a user and assumes success
func (ok AssumeSuccess) CloseAllSessions(
	user store.ID,
) []string {
	return ok.h.closeAllSessions("", user)
}

// CloseAllSessions helps closing all sessions of a user
func (notOk AssumeFailure) CloseAllSessions(
	expectedErrorCode errors.Code,
	user store.ID,
) {
	notOk.checkErrCode(expectedErrorCode)
	notOk.h.closeAllSessions(expectedErrorCode, user)
}
