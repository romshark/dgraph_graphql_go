package setup

import (
	"github.com/romshark/dgraph_graphql_go/store/errors"
)

func (h Helper) closeSession(
	expectedErrorCode errors.Code,
	key string,
) bool {
	t := h.c.t

	var result struct {
		CloseSession bool `json:"closeSession"`
	}
	checkErr(t, expectedErrorCode, h.c.QueryVar(
		`mutation (
			$key: String!
		) {
			closeSession(
				key: $key
			)
		}`,
		map[string]interface{}{
			"key": key,
		},
		&result,
	))

	if expectedErrorCode != "" {
		return false
	}

	return result.CloseSession
}

// CloseSession helps closing a session and assumes success
func (ok AssumeSuccess) CloseSession(
	key string,
) bool {
	return ok.h.closeSession("", key)
}

// CloseSession helps closing a session
func (notOk AssumeFailure) CloseSession(
	expectedErrorCode errors.Code,
	key string,
) {
	notOk.checkErrCode(expectedErrorCode)
	notOk.h.closeSession(expectedErrorCode, key)
}
